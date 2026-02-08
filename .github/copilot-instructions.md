# Copilot Instructions for terraform-linux-provider

## Project Overview
This is a **Terraform Provider for Bare-Metal Linux** that enables infrastructure-as-code management of Linux systems. The provider interfaces with systemd and ZFS through D-Bus to manage system resources.

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
│   └── internal/      # Provider resources, data sources, and tests
├── server/            # HTTP server backend for D-Bus/systemd interaction
├── common/            # Shared utilities (ZFS client, HTTP client, marshaling)
├── go.mod            # Go module dependencies
└── README.md
```

## Building and Testing

### Build Commands
```bash
# Build the provider
go build -o terraform-provider-linux

# Install dependencies
go mod download
go mod tidy
```

### Test Commands
```bash
# Run acceptance tests (from provider directory)
cd provider
make testacc

# Run acceptance tests with specific filter
cd provider
TF_ACC=1 go test ./... -v -run TestAccZPoolDataSource -timeout 120m

# Run unit tests
go test ./...
```

### Development Workflow
1. Make code changes in appropriate directory (provider/, server/, or common/)
2. Run `go mod tidy` if dependencies changed
3. Test changes with `make testacc` from provider directory
4. Verify no compilation errors: `go build ./...`

## Code Style and Conventions

### Go Code Guidelines
- Follow standard Go conventions and use `gofmt` for formatting
- Use structured logging with zerolog, not fmt.Print statements
- Implement proper error handling with descriptive error messages
- Use the Terraform Plugin Framework patterns for resources and data sources
- Keep provider-specific code in `provider/internal/provider/`
- Place shared utilities in `common/` package

### Terraform Provider Patterns
- Implement `schema.DataSource` interface for read-only data sources (e.g., zpools)
- Implement `schema.Resource` interface for managed resources
- Use `terraform-plugin-framework` validators for input validation
- Follow Terraform naming conventions: snake_case for attributes
- Add acceptance tests for all resources and data sources in `*_test.go` files

### Testing Best Practices
- Write acceptance tests with `TF_ACC=1` environment variable
- Use `terraform-plugin-testing` framework for test cases
- Include both positive and negative test cases
- Test configurations should be in HCL format
- Acceptance tests require 120m timeout minimum

## Boundaries and Constraints

### What NOT to do:
- **Never commit secrets or credentials** to the repository
- **Do not modify `go.mod`** manually; use `go mod` commands instead
- **Do not change the provider namespace** (`terraform.nickrobison.com/nickrobison/linux`)
- **Avoid breaking changes** to existing resource schemas without versioning
- **Do not add dependencies** without considering security and compatibility
- **Never remove or modify acceptance tests** without providing equivalent coverage
- **Do not bypass D-Bus integration** for direct system calls (use the established patterns)

### Security Considerations
- Provider operates with system-level permissions (systemd/D-Bus access)
- Validate all user inputs thoroughly using framework validators
- Sanitize paths and commands before system interaction
- Log sensitive operations appropriately (avoid logging credentials)

## Key Resources and Data Sources

### Currently Implemented
- **`linux_zpool` data source**: Read ZFS storage pool information
- **`linux_zpool` resource**: Create and manage ZFS storage pools

### Adding New Resources
1. Create resource schema in `provider/internal/provider/`
2. Implement CRUD operations (Create, Read, Update, Delete)
3. Add corresponding D-Bus/systemd integration in `server/` if needed
4. Create acceptance tests in `*_test.go` file
5. Update documentation

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

## CI/CD Notes
- Currently no GitHub Actions workflows configured
- Tests are run locally via Makefile
- Consider running `make testacc` before committing changes
- Ensure Go 1.22.4 compatibility for all changes
