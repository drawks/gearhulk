# GitHub Copilot Chat Instructions for Gearhulk

You are assisting with development of Gearhulk, a modern Gearman implementation in Go. Please follow these guidelines:

## Project Context
- This is a job queue system implementing the Gearman protocol
- Written in Go with focus on Kubernetes deployment
- Uses LevelDB for persistence and includes Prometheus metrics
- Supports both synchronous and asynchronous job processing

## Code Quality Standards
1. Always suggest adding tests for new functionality
2. Follow Go best practices and idioms
3. Use proper error handling (never ignore errors)
4. Include comprehensive documentation for public APIs
5. Follow the existing code style and patterns

## CLI Standards
- Use cobra for command-line interfaces
- Provide both short (-f) and long (--flag) options
- Include usage examples in help text
- Provide sensible defaults or show help when no args given

## API Design
- Keep interfaces simple and focused
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