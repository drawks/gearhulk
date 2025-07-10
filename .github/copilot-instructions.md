# GitHub Copilot Instructions for Gearhulk

## Overview
Gearhulk is a modern implementation of Gearman in Go, designed for use in lightweight deployments.
It uses Prometheus/OpenMetrics for telemetry exposition.
It has scheduled job support via cron expressions.
It implements a usable server, client and various utilities with a single static commandline executable.
The client, server, and other reusable components are also usable as importable go libraries for direct integration.

## Code Standards and Guidelines

### 1. Test Coverage Requirements
- **Always include test coverage** that demonstrates the feature or fix which has been implemented
- Every public function, method, and interface should have corresponding tests
- Tests should cover both success and failure scenarios
- Integration tests should be included for complex features
- Test files should follow Go conventions: `*_test.go` in the same package
- Use `go test -cover` to verify coverage levels

### 2. Public Interface Documentation
- **All public interfaces must be well-documented and unsurprising**
- Use Go doc comments for all exported functions, types, and constants
- Comments should explain the purpose, parameters, return values, and any side effects
- Include usage examples in documentation where appropriate
- Follow Go documentation conventions: comments start with the name of the item being documented

### 3. Command-Line Tool Standards
When implementing command-line tools:

#### 3.a Required Features:
- **Fully documented usage statements** with clear descriptions
- **GNU-style short and long options** (e.g., `-v, --verbose`)
- **Sane default behavior** OR direct short-circuit to usage statement if no sane default exists
- Consistent flag naming across commands
- Proper error handling with actionable error messages

#### 3.b CLI Implementation Guidelines:
- Use cobra/pflag for command-line parsing
- Provide both short (`-h`) and long (`--help`) forms for all flags
- Include examples in command help text
- Use appropriate default values that work for most use cases
- Validate input parameters and provide clear error messages
- Support configuration files where appropriate

### 4. Programmatic API Standards
When designing programmatic APIs:

#### 4.a Complexity Guidelines:
- **APIs should be no more complicated than necessary** to satisfy 90% of use cases
- Provide simple interfaces for common operations
- Use composition over inheritance
- Minimize the number of required parameters
- Provide sensible defaults

#### 4.b Interface Design:
- **APIs should be obvious and well-documented**
- Use clear, descriptive names for functions and types
- Implement standard Go interfaces where possible (io.Reader, io.Writer, etc.)
- Make interfaces flexible and usable in multiple contexts
- Use context.Context for cancellation and timeouts where appropriate
- Return errors as the last value in multi-return functions

### 5. Network API Standards
When implementing network APIs:

#### 5.a Protocol Compliance:
- **Follow all applicable standards** and strictly adhere to compliance with relevant protocol standards
- Implement the Gearman protocol correctly
- Use appropriate HTTP status codes and methods
- Follow REST principles for HTTP APIs
- Implement proper content negotiation

#### 5.b Error Handling:
- **Produce robust and actionable error output**
- Include specific error codes and messages
- Log errors appropriately for debugging
- Provide context about what failed and why
- Use structured logging where possible
- Handle network timeouts and connection failures gracefully

### 6. Go-Specific Guidelines
- Follow Go best practices and idioms
- Use `go fmt` for code formatting
- Run `go vet` and `golint` for code quality
- Use `go mod` for dependency management
- Implement proper error handling (don't ignore errors)
- Use channels and goroutines appropriately for concurrency
- Follow the Go naming conventions
- Use interfaces to define behavior, not data

### 7. Repository-Specific Guidelines
- This is a Gearman implementation, so maintain compatibility with the Gearman protocol
- Include Prometheus/OpenMetrics telemetry for monitoring
- Support both client and worker modes
- Provide admin interface for job management
- Use LevelDB for persistent storage by default

### 8. Testing Guidelines
- Write tests that demonstrate the feature working correctly
- Include both unit tests and integration tests
- Test error conditions and edge cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies in unit tests
- ** ALL PROJECT CODE** should consider testability as a first order concern. Where possible make dependencies injectable to maximize the ease with which units may be isolated.
- Provide test utilities for common setup/teardown

### 9. Documentation Guidelines
- Update README.md when adding new features
- Include code examples in documentation
- Document configuration options and their effects
- Provide troubleshooting information
- Keep documentation up-to-date with code changes

## Example Code Patterns

### CLI Command Example:
```go
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Start the Gearman server",
    Long: `Start the Gearman server with the specified configuration.

The server will listen for job submissions from clients and dispatch
them to available workers. It includes a web interface for monitoring
and managing jobs.

Examples:
  # Start server with default settings
  gearhulk server

  # Start server on specific address
  gearhulk server --addr 0.0.0.0:4730

  # Start server with custom storage directory
  gearhulk server --storage-dir /var/lib/gearhulk`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

### API Interface Example:
```go
// JobClient represents a client for submitting jobs to the Gearman server.
// It provides methods for submitting jobs and receiving responses.
type JobClient interface {
    // Submit submits a job to the server and returns a job handle.
    // The job will be executed by an available worker.
    Submit(ctx context.Context, funcName string, data []byte) (JobHandle, error)
    
    // SubmitBackground submits a job to run in the background.
    // Returns immediately with a job handle for tracking.
    SubmitBackground(ctx context.Context, funcName string, data []byte) (JobHandle, error)
    
    // Status returns the current status of a job.
    Status(ctx context.Context, handle JobHandle) (JobStatus, error)
    
    // Close closes the client connection.
    Close() error
}
```

These guidelines ensure that all code contributions maintain high quality, consistency, and usability standards for the Gearhulk project.
