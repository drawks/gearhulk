/*
Copyright © 2024 Dave Rawks <dave@rawks.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	rt "github.com/drawks/gearhulk/pkg/runtime"
	"github.com/drawks/gearhulk/worker"
	"github.com/spf13/cobra"
)

type workerConfig struct {
	ServerAddr string
	EofMode    bool
}

var workerCfg workerConfig

var workerCmd = &cobra.Command{
	Use:   "worker [flags] <worker-name> <command>",
	Short: "Create a worker that executes shell commands for Gearman jobs",
	Long: `Creates a Gearman worker that executes shell commands for each job.
The worker accepts job data as line-based records and passes them to the 
specified command's stdin, returning the command's stdout as the job result.

Examples:
  gearhulk worker square '/usr/bin/awk { print int($1)*int($1) }'
  gearhulk worker reverse '/usr/bin/rev'
  gearhulk worker --eof count '/usr/bin/wc -c'`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		workerName := args[0]
		command := args[1]
		
		log.Printf("Starting worker '%s' with command: %s", workerName, command)
		if workerCfg.EofMode {
			log.Println("EOF mode enabled - new subprocess for each job")
		}
		
		// Create worker
		w := worker.New(worker.Unlimited)
		defer w.Close()
		
		// Add server
		if err := w.AddServer(rt.Network, workerCfg.ServerAddr); err != nil {
			log.Fatalf("Failed to add server: %v", err)
		}
		
		// Create job handler
		jobHandler := createJobHandler(command, workerCfg.EofMode)
		
		// Add function to worker
		if err := w.AddFunc(workerName, jobHandler, 0); err != nil {
			log.Fatalf("Failed to add function: %v", err)
		}
		
		// Error handler
		w.ErrorHandler = func(e error) {
			log.Printf("Worker error: %v", e)
		}
		
		// Ready
		if err := w.Ready(); err != nil {
			log.Fatalf("Failed to ready worker: %v", err)
		}
		
		// Handle graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		
		// Start worker
		go w.Work()
		log.Printf("Worker '%s' started and ready for jobs", workerName)
		
		// Wait for shutdown signal
		<-c
		log.Println("Shutting down worker...")
		w.Close()
	},
}

func createJobHandler(command string, eofMode bool) func(worker.Job) ([]byte, error) {
	var cmdMutex sync.Mutex
	var persistentCmd *exec.Cmd
	var persistentStdin io.WriteCloser
	var persistentStdout io.ReadCloser
	
	// For non-EOF mode, start a persistent subprocess
	if !eofMode {
		var err error
		persistentCmd, persistentStdin, persistentStdout, err = startSubprocess(command)
		if err != nil {
			log.Printf("Failed to start persistent subprocess: %v", err)
			return func(job worker.Job) ([]byte, error) {
				return nil, fmt.Errorf("failed to initialize persistent subprocess: %v", err)
			}
		}
	}
	
	return func(job worker.Job) ([]byte, error) {
		cmdMutex.Lock()
		defer cmdMutex.Unlock()
		
		data := string(job.Data())
		
		if eofMode {
			// EOF mode: create new subprocess for each job
			return processJobWithNewSubprocess(command, data)
		} else {
			// Persistent mode: use existing subprocess
			return processJobWithPersistentSubprocess(persistentCmd, persistentStdin, persistentStdout, data)
		}
	}
}

func startSubprocess(command string) (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	// Use shell to execute command for better parsing
	cmd := exec.Command("sh", "-c", command)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, nil, nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return nil, nil, nil, fmt.Errorf("failed to start command: %v", err)
	}
	
	return cmd, stdin, stdout, nil
}

func processJobWithNewSubprocess(command string, data string) ([]byte, error) {
	// Use shell to execute command for better parsing
	cmd := exec.Command("sh", "-c", command)
	
	// Set up pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, fmt.Errorf("failed to start command: %v", err)
	}
	
	// Send data to stdin and close (EOF)
	if _, err := stdin.Write([]byte(data)); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, fmt.Errorf("failed to write to stdin: %v", err)
	}
	stdin.Close()
	
	// Read output
	outputBytes, err := io.ReadAll(stdout)
	if err != nil {
		stdout.Close()
		stderr.Close()
		return nil, fmt.Errorf("failed to read stdout: %v", err)
	}
	
	// Read stderr for error reporting
	stderrBytes, _ := io.ReadAll(stderr)
	
	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		if len(stderrBytes) > 0 {
			return nil, fmt.Errorf("command failed: %v, stderr: %s", err, string(stderrBytes))
		}
		return nil, fmt.Errorf("command failed: %v", err)
	}
	
	return outputBytes, nil
}

func processJobWithPersistentSubprocess(cmd *exec.Cmd, stdin io.WriteCloser, stdout io.ReadCloser, data string) ([]byte, error) {
	// Check if subprocess is still running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		return nil, fmt.Errorf("subprocess has exited")
	}
	
	// For persistent mode, send the entire job data and read the response
	// This works well with commands like 'rev' that process line by line
	// or commands that expect a full input and produce output
	
	// Send data to subprocess
	if _, err := stdin.Write([]byte(data)); err != nil {
		return nil, fmt.Errorf("failed to write to stdin: %v", err)
	}
	
	// Add a newline if the data doesn't end with one
	if !strings.HasSuffix(data, "\n") {
		if _, err := stdin.Write([]byte("\n")); err != nil {
			return nil, fmt.Errorf("failed to write newline to stdin: %v", err)
		}
	}
	
	// Read the response line by line until we get the expected number of lines
	inputLines := strings.Split(strings.TrimSpace(data), "\n")
	var results []string
	
	reader := bufio.NewReader(stdout)
	for i := 0; i < len(inputLines); i++ {
		if inputLines[i] == "" {
			continue
		}
		
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdout: %v", err)
		}
		
		// Remove trailing newline
		line = strings.TrimSuffix(line, "\n")
		results = append(results, line)
	}
	
	return []byte(strings.Join(results, "\n")), nil
}

func init() {
	rootCmd.AddCommand(workerCmd)
	
	workerCmd.Flags().StringVar(&workerCfg.ServerAddr, "server", "127.0.0.1:4730", "Gearman server address")
	workerCmd.Flags().BoolVar(&workerCfg.EofMode, "eof", false, "Terminate each record with EOF and fork new subprocess")
}