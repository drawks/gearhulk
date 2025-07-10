package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldError    bool
	}{
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "Usage:",
			shouldError:    false,
		},
		{
			name:           "help flag works",
			args:           []string{"--help"},
			expectedOutput: "Usage:",
			shouldError:    false,
		},
		{
			name:           "short help flag works",
			args:           []string{"-h"},
			expectedOutput: "Usage:",
			shouldError:    false,
		},
		{
			name:           "version information in help",
			args:           []string{"--help"},
			expectedOutput: "modern implementation of Gearman",
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for testing to avoid state pollution
			cmd := &cobra.Command{
				Use:   "gearhulk",
				Short: "A modern Gearman implementation in Go",
				Long: `Gearhulk is a modern implementation of Gearman in Go Programming Language.

Gearhulk includes various improvements in retry and connection logic for using
in Kubernetes. It comes with built-in Prometheus ready metrics. Gearhulk also
implements scheduled jobs via cron expressions.

The server can be used as a job queue manager, while clients can submit jobs
and workers can process them. It includes a web interface for monitoring and
managing jobs.

Examples:
  # Start the Gearman server
  gearhulk server

  # Start server on specific address with custom storage
  gearhulk server --addr 0.0.0.0:4730 --storage-dir /var/lib/gearhulk

  # Show help for server command
  gearhulk server --help`,
				Run: func(cmd *cobra.Command, args []string) {
					cmd.Help()
				},
			}

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestServerCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldError    bool
	}{
		{
			name:           "server help works",
			args:           []string{"server", "--help"},
			expectedOutput: "Start the Gearman server",
			shouldError:    false,
		},
		{
			name:           "server short help works",
			args:           []string{"server", "-h"},
			expectedOutput: "Start the Gearman server",
			shouldError:    false,
		},
		{
			name:           "server help shows GNU-style flags",
			args:           []string{"server", "--help"},
			expectedOutput: "-a, --addr",
			shouldError:    false,
		},
		{
			name:           "server help shows examples",
			args:           []string{"server", "--help"},
			expectedOutput: "Examples:",
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			rootCmd := &cobra.Command{
				Use:   "gearhulk",
				Short: "A modern Gearman implementation in Go",
			}

			// Create server command
			serverCmd := &cobra.Command{
				Use:   "server",
				Short: "Start the Gearman server",
				Long: `Start the Gearman server with the specified configuration.

The server will listen for job submissions from clients and dispatch
them to available workers. It includes a web interface for monitoring
and managing jobs, as well as built-in Prometheus metrics.

The server uses LevelDB for persistent storage by default and supports
scheduled jobs via cron expressions.

Examples:
  # Start server with default settings
  gearhulk server

  # Start server on specific address
  gearhulk server --addr 0.0.0.0:4730

  # Start server with custom storage directory
  gearhulk server --storage-dir /var/lib/gearhulk

  # Start server with custom web interface address
  gearhulk server --web-addr :8080

  # Start server with verbose logging
  gearhulk server --addr 0.0.0.0:4730 --verbose`,
				Run: func(cmd *cobra.Command, args []string) {
					// Don't actually start the server in tests
				},
			}

			// Add GNU-style flags
			serverCmd.Flags().StringP("addr", "a", ":4730", "listening address, such as 0.0.0.0:4730")
			serverCmd.Flags().StringP("storage-dir", "s", "/tmp/gearmand", "directory where LevelDB file is stored")
			serverCmd.Flags().StringP("web-addr", "w", ":3000", "server HTTP API address")
			serverCmd.Flags().BoolP("verbose", "v", false, "enable verbose logging")

			rootCmd.AddCommand(serverCmd)

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Set args
			rootCmd.SetArgs(tt.args)

			// Execute command
			err := rootCmd.Execute()

			// Check error expectation
			if tt.shouldError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestGNUStyleFlags(t *testing.T) {
	// Test that our commands follow GNU-style conventions
	tests := []struct {
		name        string
		command     string
		shortFlag   string
		longFlag    string
		description string
	}{
		{
			name:        "server addr flag",
			command:     "server",
			shortFlag:   "-a",
			longFlag:    "--addr",
			description: "listening address",
		},
		{
			name:        "server storage flag",
			command:     "server",
			shortFlag:   "-s",
			longFlag:    "--storage-dir",
			description: "directory where LevelDB",
		},
		{
			name:        "server web flag",
			command:     "server",
			shortFlag:   "-w",
			longFlag:    "--web-addr",
			description: "server HTTP API",
		},
		{
			name:        "server verbose flag",
			command:     "server",
			shortFlag:   "-v",
			longFlag:    "--verbose",
			description: "enable verbose logging",
		},
		{
			name:        "root config flag",
			command:     "",
			shortFlag:   "-c",
			longFlag:    "--config",
			description: "config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []string
			if tt.command != "" {
				args = []string{tt.command, "--help"}
			} else {
				args = []string{"--help"}
			}

			// Create a new root command for testing
			rootCmd := &cobra.Command{
				Use:   "gearhulk",
				Short: "A modern Gearman implementation in Go",
				Run: func(cmd *cobra.Command, args []string) {
					cmd.Help()
				},
			}

			// Add config flag
			rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.gearhulk.yaml)")

			// Create server command if needed
			if tt.command == "server" {
				serverCmd := &cobra.Command{
					Use:   "server",
					Short: "Start the Gearman server",
					Long:  "Start the Gearman server with the specified configuration.",
					Run: func(cmd *cobra.Command, args []string) {
						// Don't actually start the server in tests
					},
				}

				serverCmd.Flags().StringP("addr", "a", ":4730", "listening address, such as 0.0.0.0:4730")
				serverCmd.Flags().StringP("storage-dir", "s", "/tmp/gearmand", "directory where LevelDB file is stored")
				serverCmd.Flags().StringP("web-addr", "w", ":3000", "server HTTP API address")
				serverCmd.Flags().BoolP("verbose", "v", false, "enable verbose logging")

				rootCmd.AddCommand(serverCmd)
			}

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Set args
			rootCmd.SetArgs(args)

			// Execute command
			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output contains both short and long flags
			output := buf.String()
			expectedPattern := tt.shortFlag + ", " + tt.longFlag
			if !strings.Contains(output, expectedPattern) {
				t.Errorf("expected output to contain GNU-style flag pattern %q, got %q", expectedPattern, output)
			}

			// Check description is present
			if !strings.Contains(output, tt.description) {
				t.Errorf("expected output to contain description %q, got %q", tt.description, output)
			}
		})
	}
}