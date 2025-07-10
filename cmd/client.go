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
	"os"
	"time"

	"github.com/drawks/gearhulk/client"
	rt "github.com/drawks/gearhulk/pkg/runtime"
	"github.com/spf13/cobra"
)

var (
	serverAddr string
	delimiter  string
	timeout    time.Duration
	priority   string
)

type jobResult struct {
	data []byte
	err  error
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client [function]",
	Short: "Submit jobs to gearman servers via command line",
	Long: `A basic command line client for submitting jobs to gearman servers.
	
Reads job data from stdin (line-delimited by default) and submits each line
as a job to the specified function. Results are output to stdout.

Example:
  seq 1 10 | gearhulk client square
  echo "hello world" | gearhulk client reverse`,
	Args: cobra.ExactArgs(1),
	RunE: runClient,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVar(&serverAddr, "server", "127.0.0.1:4730", "gearman server address")
	clientCmd.Flags().StringVar(&delimiter, "delimiter", "\n", "input delimiter (default: newline)")
	clientCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "job timeout")
	clientCmd.Flags().StringVar(&priority, "priority", "normal", "job priority (low, normal, high)")
}

func runClient(cmd *cobra.Command, args []string) error {
	funcName := args[0]

	// Create client connection
	c, err := client.New("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server %s: %w", serverAddr, err)
	}
	defer c.Close()

	// Set response timeout
	c.ResponseTimeout = timeout

	// Set error handler
	c.ErrorHandler = func(e error) {
		fmt.Fprintf(os.Stderr, "Client error: %v\n", e)
	}

	// Determine job priority
	var jobPriority byte
	switch priority {
	case "low":
		jobPriority = rt.JobLow
	case "high":
		jobPriority = rt.JobHigh
	default:
		jobPriority = rt.JobNormal
	}

	// Read from stdin and submit jobs
	scanner := bufio.NewScanner(os.Stdin)
	
	// Set custom delimiter if specified
	if delimiter != "\n" {
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := len([]byte(delimiter)); i > 0 {
				for j := 0; j <= len(data)-i; j++ {
					if string(data[j:j+i]) == delimiter {
						return j + i, data[0:j], nil
					}
				}
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		})
	}

	// Process each line from stdin sequentially
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Channel to receive the result for this specific job
		resultChan := make(chan jobResult, 1)
		
		// Create response handler
		responseHandler := func(resp *client.Response) {
			switch resp.DataType {
			case rt.PT_WorkComplete:
				resultChan <- jobResult{data: resp.Data, err: nil}
			case rt.PT_WorkFail:
				resultChan <- jobResult{data: nil, err: fmt.Errorf("job failed")}
			case rt.PT_WorkException:
				resultChan <- jobResult{data: nil, err: fmt.Errorf("job exception: %s", string(resp.Data))}
			}
		}

		// Submit job
		_, err := c.Do(funcName, []byte(line), jobPriority, responseHandler)
		if err != nil {
			return fmt.Errorf("error submitting job: %w", err)
		}

		// Wait for result
		select {
		case result := <-resultChan:
			if result.err != nil {
				fmt.Fprintf(os.Stderr, "Job error: %v\n", result.err)
			} else {
				fmt.Println(string(result.data))
			}
		case <-time.After(timeout):
			fmt.Fprintf(os.Stderr, "Job timeout after %v\n", timeout)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}

	return nil
}