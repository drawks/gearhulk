# GitHub Copilot Chat Instructions for Gearhulk

You are assisting with development of Gearhulk, a modern Gearman implementation in Go. Please follow these guidelines:

## Project Context
- This is a job queue system implementing the Gearman protocol
- Written in Go with a focus on lightweight deployments
- Uses LevelDB for persistence and includes Prometheus metrics
- Supports both synchronous and asynchronous job processing

## Code Quality Standards
1. All new functionality requires test coverage
2. All bug fixes require regression tests that demonstrate the bug
3. Follow Go best practices and idioms
4. Use proper error handling (never ignore errors)
5. Include comprehensive documentation for public APIs
6. Do not be overly concerned with maintaining existing codestyle if/when it significantly deviates from modern acceptable standards
7. Do not mix functional changes with pure formatting changes in the same commit

## CLI Standards
- Use cobra for command-line interfaces
- Provide both short (-f) and long (--flag) options
- Include usage examples in help text
- Provide sensible defaults or show help when no args given

## API Design
- Keep interfaces simple and focused, if possible and sensible always strive for compatibility with an existing common interface
- Use context.Context for cancellation/timeouts
- Follow Go naming conventions
- Make APIs flexible and composable

## Testing
- Write table-driven tests where appropriate
- Test both success and failure cases
- Use testify for assertions
- Include integration tests for complex features

## Documentation
- Use Go doc comments for all exported items
- Include examples in documentation
- Update README.md for new features
- Document configuration options

When suggesting code changes, always consider:
- Is this the minimal change needed?
- Does this follow the project's patterns?
- Are there tests for this functionality?
- Is the API easy to use and understand?
