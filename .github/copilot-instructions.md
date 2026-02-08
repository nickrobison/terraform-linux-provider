# Copilot Instructions for terraform-linux-provider

## Project Overview
This is a **Terraform Provider for Bare-Metal Linux** that enables infrastructure-as-code management of Linux systems. The project consists of two main components:

1. **Linux Server** (`server/`): A Go-based HTTP server that interfaces with Linux system services via D-Bus
2. **Terraform Provider** (`provider/`): A Terraform plugin that communicates with the Linux server to manage infrastructure

The server provides a REST API layer over D-Bus operations for systemd and ZFS management, while the provider translates Terraform configurations into API calls to the server.

## Tech Stack
- **Language**: Go 1.22.4
- **Framework**: HashiCorp Terraform Plugin Framework v1.10.0
- **System Integration**: D-Bus (godbus v5.1.0) for systemd communication
- **Storage**: ZFS pool management (zpools data source and resource)
- **Logging**: Zerolog v1.33.0 for structured logging
- **Testing**: Terraform Plugin Testing Framework v1.9.0

## Project Structure
```
terraform-linux-provider/
├── provider/          # Terraform provider implementation
│   └── internal/      # Provider resources, data sources, and acceptance tests
├── server/            # Go-based Linux server - HTTP API for D-Bus/systemd
│   ├── main.go        # Server entry point
│   ├── routes.go      # HTTP route definitions
│   ├── middleware/    # HTTP middleware (logging, etc.)
│   ├── zfs/           # ZFS pool management via D-Bus
│   └── bus/           # D-Bus utility functions
├── common/            # Shared utilities (ZFS client, HTTP client, marshaling)
├── go.mod            # Go module dependencies
└── README.md
```

## Building and Testing

### Build Commands
```bash
# Build the Terraform provider
cd provider
go build -o terraform-provider-linux

# Build the Linux server
cd server
go build -o linux-server

# Install dependencies (from repository root)
go mod download
go mod tidy

# Verify all packages compile
go build ./...
```

### Test Commands
```bash
# Run all acceptance tests (from provider directory)
cd provider
make testacc

# Run specific acceptance test
cd provider
TF_ACC=1 go test ./... -v -run TestAccZPoolDataSource -timeout 120m

# Run unit tests for all packages
go test ./...

# Run tests for specific package
go test ./server/zfs
go test ./common
go test ./provider/internal/provider

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
```

### Development Workflow
1. Make code changes in appropriate directory (provider/, server/, or common/)
2. Run `go mod tidy` if dependencies changed
3. Write or update tests for your changes (see Testing Requirements below)
4. Run unit tests: `go test ./path/to/package`
5. For provider changes, run acceptance tests: `cd provider && make testacc`
6. Verify no compilation errors: `go build ./...`
7. Ensure all tests pass before submitting changes

## Architecture

### Component Interaction
```
Terraform Config → Provider (provider/) → HTTP Client (common/client.go) 
                                                ↓
                                         Server (server/)
                                                ↓
                                         D-Bus Client
                                                ↓
                                    Linux System (systemd/ZFS)
```

The provider never directly interacts with Linux system services. All system interactions are mediated through the server component, which provides a clean separation of concerns and enables better testing.

## Code Standards

### Required Before Each Commit
- Run `go test ./... -short -v -race` to ensure all tests pass (matches CI)
- For provider changes, run `cd provider && make testacc` to verify acceptance tests
- Run `go build ./...` to verify no compilation errors
- Run `go mod tidy` if you added or removed dependencies
- Ensure new code has appropriate test coverage
- CI will automatically run unit tests (with race detection) and integration tests on PRs

## Code Style and Conventions

### Go Code Guidelines
- Follow standard Go conventions and use `gofmt` for formatting
- Use structured logging with zerolog, not fmt.Print statements
- Implement proper error handling with descriptive error messages
- Keep provider-specific code in `provider/internal/provider/`
- Keep server-specific code in `server/` package
- Place shared utilities in `common/` package
- Use dependency injection for testability (see server/main.go)

### Server Component Guidelines
- The server runs as an HTTP API on localhost:8080
- All D-Bus interactions happen in the server, not the provider
- Use zerolog for structured logging with appropriate log levels
- Implement proper error handling and HTTP status codes
- Follow RESTful API patterns for routes
- Keep handlers focused and testable
- D-Bus clients should be mockable for testing

### Terraform Provider Patterns
- The provider communicates with the server via HTTP (common/client.go)
- Implement `schema.DataSource` interface for read-only data sources (e.g., zpools)
- Implement `schema.Resource` interface for managed resources
- Use `terraform-plugin-framework` validators for input validation
- Follow Terraform naming conventions: snake_case for attributes
- Add acceptance tests for all resources and data sources in `*_test.go` files
- Provider should not contain D-Bus logic - delegate to server

### Testing Requirements
**Test coverage is critical in this project. All code changes must include appropriate tests.**

#### Acceptance Tests (Provider)
- Write acceptance tests with `TF_ACC=1` environment variable
- Use `terraform-plugin-testing` framework for test cases
- Include both positive and negative test cases
- Test configurations should be in HCL format
- Acceptance tests require 120m timeout minimum
- Test file naming: `*_test.go` in same package as implementation
- Example: `provider/internal/provider/zpool_data_source_test.go`

#### Unit Tests (Server and Common)
- Write unit tests for all server endpoints and handlers
- Test D-Bus client interactions with mock connections where possible
- Test common utilities (marshaling, HTTP clients) thoroughly
- Use table-driven tests for multiple test cases
- Test file naming: `*_test.go` in same package as implementation
- Aim for high coverage of business logic and error paths

#### Test Guidelines
- Always add tests when adding new features
- Update tests when modifying existing functionality
- Never remove or modify existing tests without equivalent coverage
- Test both success and failure scenarios
- Validate error messages and error handling
- Use descriptive test names that explain what is being tested

## Boundaries and Constraints

### What NOT to do:
- **Never commit secrets or credentials** to the repository
- **Do not modify `go.mod`** manually; use `go mod` commands instead
- **Do not change the provider namespace** (`terraform.nickrobison.com/nickrobison/linux`)
- **Avoid breaking changes** to existing resource schemas without versioning
- **Do not add dependencies** without considering security and compatibility
- **Never remove or modify existing tests** without providing equivalent coverage
- **Do not bypass the server layer** - provider should not directly call D-Bus
- **Do not add D-Bus code to the provider** - all D-Bus interactions belong in server/
- **Do not skip writing tests** - test coverage is mandatory for all changes

### Security Considerations
- The server operates with system-level permissions (systemd/D-Bus access)
- Validate all user inputs in both provider and server
- Use terraform-plugin-framework validators in the provider
- Sanitize paths and commands in the server before system interaction
- Log sensitive operations appropriately (avoid logging credentials)
- Server should validate all incoming HTTP requests
- Provider should use HTTPS in production (currently localhost for development)

## Key Resources and Data Sources

### Currently Implemented
- **`linux_zpool` data source**: Read ZFS storage pool information
- **`linux_zpool` resource**: Create and manage ZFS storage pools

### Adding New Resources or Features
1. **Server Side** (if new D-Bus integration needed):
   - Add D-Bus client in `server/<component>/dbus_client.go`
   - Add HTTP handler in `server/<component>/handler.go`
   - Add route in `server/routes.go`
   - Write unit tests for handler and D-Bus client
   
2. **Provider Side**:
   - Create resource schema in `provider/internal/provider/`
   - Implement CRUD operations (Create, Read, Update, Delete)
   - Use HTTP client from `common/client.go` to call server API
   - Create acceptance tests in `*_test.go` file
   
3. **Common Utilities** (if shared code needed):
   - Add shared types or utilities in `common/`
   - Write unit tests for common utilities
   
4. **Documentation**:
   - Update README.md with usage examples
   - Document resource attributes and behavior

## Documentation
- Keep README.md updated with provider usage examples
- Document all resource attributes and data source outputs
- Include example Terraform configurations for common use cases
- Document any system requirements (Linux version, ZFS, systemd)

## Dependencies and Updates
- Use `go get -u` cautiously; test thoroughly after updates
- Keep Terraform Plugin Framework synchronized with Terraform versions
- Ensure D-Bus library compatibility with target Linux distributions
- Monitor HashiCorp provider best practices and migration guides

## CI/CD and Quality Checks
- **GitHub Actions workflows** are configured for automated testing:
  - **Unit Tests** (`.github/workflows/unit-tests.yml`): Runs on Go 1.22.x and 1.23.x, includes race detection and coverage reporting
  - **Integration Tests** (`.github/workflows/integration-tests.yml`): Runs acceptance tests with `TF_ACC=1` on Go 1.22.x
- Workflows run automatically on:
  - Push to `main` branch
  - Pull request creation and updates
- Local testing before committing:
  - Run unit tests: `go test ./... -short -v -race` (matches CI unit tests)
  - Run acceptance tests: `cd provider && make testacc` (requires Linux with ZFS)
  - Run all tests: `go test ./...` from repository root
- Ensure Go 1.22.4 or later compatibility (go.mod requires 1.22.4; CI tests against 1.22.x and 1.23.x)
- Acceptance tests require a Linux environment with ZFS and D-Bus available

## Common Pitfalls to Avoid
- Adding D-Bus interactions in the provider instead of the server
- Forgetting to add tests for new functionality
- Not testing both success and error paths
- Breaking the provider/server architectural boundary
- Modifying go.mod manually instead of using `go mod` commands
- Not running acceptance tests before submitting provider changes
