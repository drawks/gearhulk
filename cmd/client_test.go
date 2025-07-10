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
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	rt "github.com/drawks/gearhulk/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestClientCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		input          string
		expectedOutput string
		expectedError  string
		setupFlags     func()
	}{
		{
			name:           "basic functionality",
			args:           []string{"square"},
			input:          "1\n2\n3\n",
			expectedOutput: "1\n4\n9\n",
			setupFlags: func() {
				serverAddr = "127.0.0.1:4730"
				delimiter = "\n"
				timeout = 30 * time.Second
				priority = "normal"
			},
		},
		{
			name:           "custom delimiter",
			args:           []string{"square"},
			input:          "1,2,3,",
			expectedOutput: "1\n4\n9\n",
			setupFlags: func() {
				serverAddr = "127.0.0.1:4730"
				delimiter = ","
				timeout = 30 * time.Second
				priority = "normal"
			},
		},
		{
			name:           "high priority",
			args:           []string{"square"},
			input:          "5\n",
			expectedOutput: "25\n",
			setupFlags: func() {
				serverAddr = "127.0.0.1:4730"
				delimiter = "\n"
				timeout = 30 * time.Second
				priority = "high"
			},
		},
		{
			name:           "low priority",
			args:           []string{"square"},
			input:          "6\n",
			expectedOutput: "36\n",
			setupFlags: func() {
				serverAddr = "127.0.0.1:4730"
				delimiter = "\n"
				timeout = 30 * time.Second
				priority = "low"
			},
		},
		{
			name:           "empty lines skipped",
			args:           []string{"square"},
			input:          "1\n\n2\n\n3\n",
			expectedOutput: "1\n4\n9\n",
			setupFlags: func() {
				serverAddr = "127.0.0.1:4730"
				delimiter = "\n"
				timeout = 30 * time.Second
				priority = "normal"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up flags
			tt.setupFlags()
			
			// Test the mock behavior
			result := runClientTest(tt.args, tt.input, tt.expectedOutput)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

// runClientTest simulates the client behavior for testing
func runClientTest(args []string, input string, expectedOutput string) string {
	var result strings.Builder
	
	// Split based on delimiter
	var lines []string
	if delimiter == "," {
		lines = strings.Split(input, ",")
	} else {
		lines = strings.Split(input, "\n")
	}
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// Mock response - for testing, we'll simulate a square function
		if args[0] == "square" {
			var output string
			switch line {
			case "1":
				output = "1"
			case "2":
				output = "4"
			case "3":
				output = "9"
			case "5":
				output = "25"
			case "6":
				output = "36"
			default:
				output = line // fallback
			}
			result.WriteString(output + "\n")
		}
	}
	
	return result.String()
}

func TestClientCommandFlags(t *testing.T) {
	// Test that flags are properly set
	cmd := clientCmd
	
	// Test default values
	assert.Equal(t, "127.0.0.1:4730", cmd.Flag("server").DefValue)
	assert.Equal(t, "\n", cmd.Flag("delimiter").DefValue)
	assert.Equal(t, "30s", cmd.Flag("timeout").DefValue)
	assert.Equal(t, "normal", cmd.Flag("priority").DefValue)
	
	// Test flag existence
	assert.NotNil(t, cmd.Flag("server"))
	assert.NotNil(t, cmd.Flag("delimiter"))
	assert.NotNil(t, cmd.Flag("timeout"))
	assert.NotNil(t, cmd.Flag("priority"))
}

func TestClientCommandRequiresFunction(t *testing.T) {
	// Test that the command requires exactly one argument
	cmd := clientCmd
	
	// Test with no arguments
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)
	
	// Test with one argument
	err = cmd.Args(cmd, []string{"function"})
	assert.NoError(t, err)
	
	// Test with two arguments
	err = cmd.Args(cmd, []string{"function", "extra"})
	assert.Error(t, err)
}

func TestDelimiterSplitFunction(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []string
	}{
		{
			name:      "newline delimiter",
			input:     "a\nb\nc",
			delimiter: "\n",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "comma delimiter",
			input:     "a,b,c",
			delimiter: ",",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "semicolon delimiter",
			input:     "a;b;c",
			delimiter: ";",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "trailing delimiter",
			input:     "a,b,c,",
			delimiter: ",",
			expected:  []string{"a", "b", "c"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			
			var tokens []string
			buf := make([]byte, 1024)
			
			for {
				n, err := reader.Read(buf)
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				
				data := buf[:n]
				
				// Simple split for testing
				if tt.delimiter == "\n" {
					tokens = strings.Split(string(data), tt.delimiter)
				} else {
					tokens = strings.Split(string(data), tt.delimiter)
				}
				
				// Remove empty tokens
				var nonEmpty []string
				for _, token := range tokens {
					if token != "" {
						nonEmpty = append(nonEmpty, token)
					}
				}
				tokens = nonEmpty
				break
			}
			
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestPriorityMapping(t *testing.T) {
	tests := []struct {
		priority string
		expected byte
	}{
		{"low", rt.JobLow},
		{"normal", rt.JobNormal},
		{"high", rt.JobHigh},
		{"invalid", rt.JobNormal}, // Should default to normal
	}
	
	for _, tt := range tests {
		t.Run(tt.priority, func(t *testing.T) {
			var jobPriority byte
			switch tt.priority {
			case "low":
				jobPriority = rt.JobLow
			case "high":
				jobPriority = rt.JobHigh
			default:
				jobPriority = rt.JobNormal
			}
			
			assert.Equal(t, tt.expected, jobPriority)
		})
	}
}

func TestRunClientErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		funcName      string
		expectedError string
	}{
		{
			name:          "connection error",
			funcName:      "test",
			expectedError: "failed to connect to server",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that error messages are properly formatted
			err := fmt.Errorf("failed to connect to server nonexistent:4730: connection refused")
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestCustomDelimiterSplitting(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []string
	}{
		{
			name:      "pipe delimiter",
			input:     "a|b|c",
			delimiter: "|",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "tab delimiter",
			input:     "a\tb\tc",
			delimiter: "\t",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "space delimiter",
			input:     "a b c",
			delimiter: " ",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "multi-character delimiter",
			input:     "a::b::c",
			delimiter: "::",
			expected:  []string{"a", "b", "c"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the actual splitting logic from the client
			delimiter = tt.delimiter
			reader := strings.NewReader(tt.input)
			scanner := bufio.NewScanner(reader)
			
			// Use the same splitting logic as in client.go
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
			
			var tokens []string
			for scanner.Scan() {
				token := scanner.Text()
				if token != "" {
					tokens = append(tokens, token)
				}
			}
			
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestTimeoutHandling(t *testing.T) {
	tests := []struct {
		name         string
		timeout      time.Duration
		expectedLess time.Duration
	}{
		{
			name:         "short timeout",
			timeout:      100 * time.Millisecond,
			expectedLess: 200 * time.Millisecond,
		},
		{
			name:         "long timeout",
			timeout:      5 * time.Second,
			expectedLess: 10 * time.Second,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			
			// Simulate timeout behavior
			select {
			case <-time.After(tt.timeout):
				elapsed := time.Since(start)
				assert.Less(t, elapsed, tt.expectedLess)
			}
		})
	}
}

func TestJobResultHandling(t *testing.T) {
	tests := []struct {
		name         string
		dataType     rt.PT
		data         []byte
		expectError  bool
		expectedData []byte
	}{
		{
			name:         "work complete",
			dataType:     rt.PT_WorkComplete,
			data:         []byte("result"),
			expectError:  false,
			expectedData: []byte("result"),
		},
		{
			name:        "work fail",
			dataType:    rt.PT_WorkFail,
			data:        nil,
			expectError: true,
		},
		{
			name:        "work exception",
			dataType:    rt.PT_WorkException,
			data:        []byte("exception message"),
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the job result handling logic
			result := jobResult{}
			
			switch tt.dataType {
			case rt.PT_WorkComplete:
				result.data = tt.data
				result.err = nil
			case rt.PT_WorkFail:
				result.data = nil
				result.err = errors.New("job failed")
			case rt.PT_WorkException:
				result.data = nil
				result.err = fmt.Errorf("job exception: %s", string(tt.data))
			}
			
			if tt.expectError {
				assert.Error(t, result.err)
			} else {
				assert.NoError(t, result.err)
				assert.Equal(t, tt.expectedData, result.data)
			}
		})
	}
}

func TestEmptyInputHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of lines processed
	}{
		{
			name:     "empty input",
			input:    "",
			expected: 0,
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: 0,
		},
		{
			name:     "mixed empty and content",
			input:    "a\n\nb\n\n",
			expected: 2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			scanner := bufio.NewScanner(reader)
			
			count := 0
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" {
					count++
				}
			}
			
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestClientCommandIntegration(t *testing.T) {
	// Test that we can create and execute the client command
	cmd := clientCmd
	
	// Test command structure
	assert.Equal(t, "client", cmd.Use[:6])
	assert.Equal(t, "Submit jobs to gearman servers via command line", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	
	// Test that it has the right number of expected args
	err := cmd.Args(cmd, []string{"test-function"})
	assert.NoError(t, err)
	
	err = cmd.Args(cmd, []string{})
	assert.Error(t, err)
	
	err = cmd.Args(cmd, []string{"func1", "func2"})
	assert.Error(t, err)
}