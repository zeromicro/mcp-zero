# Error Recovery Guide

**Version**: 1.0.0
**Date**: November 15, 2025

## Overview

This guide provides comprehensive error recovery examples for common issues encountered when using mcp-zero. Each section includes the error, cause, and solution.

## Table of Contents

1. [Service Creation Errors](#service-creation-errors)
2. [Build Errors](#build-errors)
3. [Import Path Errors](#import-path-errors)
4. [Configuration Errors](#configuration-errors)
5. [Database Connection Errors](#database-connection-errors)
6. [goctl Discovery Errors](#goctl-discovery-errors)
7. [Port Conflicts](#port-conflicts)
8. [Module Initialization Errors](#module-initialization-errors)

---

## Service Creation Errors

### Error: Invalid Service Name

**Symptom**:
```
Error: service name cannot contain hyphens, try: user_service
```

**Cause**: Service names must be valid Go identifiers (no hyphens)

**Solution**:
```
❌ Wrong: Create API service "user-service"
✅ Correct: Create API service "userservice" or "user_service"
```

**Automatic Fix**: mcp-zero suggests valid alternatives

---

### Error: Service Already Exists

**Symptom**:
```
Error: directory /path/to/userservice already exists
```

**Cause**: Trying to create a service in an existing directory

**Solution**:
1. Choose a different service name
2. Remove the existing directory
3. Specify a different output directory

```bash
# Option 1: Different name
Create API service "userservice2"

# Option 2: Different location
Create API service "userservice" in /path/to/services/
```

---

## Build Errors

### Error: Cannot Find Package

**Symptom**:
```
Error: build failed: package internal/config: cannot find package
```

**Cause**: Import paths are incorrect after generation

**Recovery**: mcp-zero automatically fixes imports, but if manual fix needed:

```bash
cd yourservice
go mod init github.com/yourname/yourservice
go mod tidy
```

**Prevention**: Always use absolute module paths

---

### Error: Undefined Symbol

**Symptom**:
```
Error: undefined: config.Config
```

**Cause**: Missing type definition or import

**Solution**:
```go
// Add missing import
import "yourservice/internal/config"

// OR define the type
type Config struct {
    Port int
}
```

**Automatic Fix**: mcp-zero runs `go mod tidy` to resolve dependencies

---

## Import Path Errors

### Error: Import Cycle Detected

**Symptom**:
```
Error: import cycle not allowed
package main
    imports yourservice/internal/handler
    imports yourservice/internal/logic
    imports yourservice/internal/handler
```

**Cause**: Circular dependencies between packages

**Solution**: Restructure code to remove circular dependencies

```go
// BAD: handler -> logic -> handler (cycle)

// GOOD: handler -> logic -> repository (no cycle)
// Extract shared code into a separate package
```

**Prevention**:
- Keep dependencies uni-directional
- Use interfaces to break cycles
- Extract shared code into common packages

---

### Error: Cannot Find Module

**Symptom**:
```
Error: go.mod file not found in current directory or any parent directory
```

**Cause**: Module not initialized

**Solution**: mcp-zero automatically runs `go mod init`, but if manual fix needed:

```bash
cd yourservice
go mod init github.com/yourname/yourservice
```

---

## Configuration Errors

### Error: Invalid YAML Syntax

**Symptom**:
```
Error: yaml: line 5: did not find expected key
```

**Cause**: Malformed YAML configuration

**Solution**: Fix YAML syntax

```yaml
# ❌ Wrong: Missing colon
Server
  Port 8080

# ✅ Correct: Proper YAML syntax
Server:
  Port: 8080
  Host: localhost
```

**Tool**: Use YAML validator or IDE with YAML support

---

### Error: Missing Required Config Field

**Symptom**:
```
Error: configuration validation failed: Port is required
```

**Cause**: Required configuration field not set

**Solution**: Add missing fields

```yaml
# Add required fields
Server:
  Port: 8080        # Required
  Host: localhost   # Optional
  Timeout: 30s      # Optional
```

**Documentation**: Check config template for required fields

---

## Database Connection Errors

### Error: Connection Refused

**Symptom**:
```
Error: failed to connect to database: connection refused
```

**Cause**: Database server not running or wrong connection string

**Solution**:

```bash
# 1. Check database is running
mysql.server status  # macOS
systemctl status mysql  # Linux

# 2. Verify connection string format
mysql://user:password@localhost:3306/database

# 3. Test connection
mysql -u user -p -h localhost database
```

**Common Issues**:
- Wrong host/port
- Incorrect credentials
- Database doesn't exist
- Firewall blocking connection

---

### Error: Unknown Database

**Symptom**:
```
Error: unknown database 'mydb'
```

**Cause**: Database doesn't exist

**Solution**:

```sql
-- Create the database first
CREATE DATABASE mydb;

-- Then run model generation
```

**Prevention**: Always create database before generating models

---

## goctl Discovery Errors

### Error: goctl Not Found

**Symptom**:
```
Error: goctl not found in PATH or standard locations
```

**Cause**: goctl not installed or not in PATH

**Solution**:

```bash
# Install goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Verify installation
which goctl
goctl --version

# If not in PATH, set GOCTL_PATH environment variable
export GOCTL_PATH=/path/to/goctl
```

**Claude Desktop**: Add to configuration

```json
{
  "mcpServers": {
    "go-zero": {
      "command": "/path/to/mcp-zero",
      "env": {
        "GOCTL_PATH": "/Users/yourname/go/bin/goctl"
      }
    }
  }
}
```

---

### Error: goctl Version Incompatible

**Symptom**:
```
Error: goctl version 1.3.0 is not compatible, requires 1.4.0+
```

**Cause**: Outdated goctl version

**Solution**:

```bash
# Update goctl to latest version
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Verify version
goctl --version
```

---

## Port Conflicts

### Error: Port Already in Use

**Symptom**:
```
Error: port 8080 is already in use
```

**Cause**: Another process is using the port

**Solution**:

```bash
# Option 1: Find and stop the process
lsof -i :8080
kill <PID>

# Option 2: Use a different port
# Ask mcp-zero to suggest an available port
```

**Automatic Fix**: mcp-zero validates ports before creation

---

### Error: Port Out of Range

**Symptom**:
```
Error: port must be between 1024 and 65535, got 80
```

**Cause**: Trying to use privileged port or invalid range

**Solution**:

```
❌ Wrong: Use port 80 (privileged)
❌ Wrong: Use port 70000 (out of range)
✅ Correct: Use port 8080 (valid range)
```

**Valid Range**: 1024-65535 (non-privileged ports)

---

## Module Initialization Errors

### Error: go.mod Already Exists

**Symptom**:
```
Error: go.mod file already exists
```

**Cause**: Trying to initialize module in a directory that already has one

**Solution**: This is usually not a problem. mcp-zero skips initialization if go.mod exists.

**Manual Fix** (if needed):

```bash
# Remove existing go.mod
rm go.mod go.sum

# Re-initialize
go mod init github.com/yourname/yourservice
go mod tidy
```

---

### Error: Invalid Module Path

**Symptom**:
```
Error: invalid module path: "my service"
```

**Cause**: Module path contains invalid characters (spaces)

**Solution**:

```
❌ Wrong: "my service"
✅ Correct: "myservice" or "my-service" or "github.com/user/my-service"
```

**Best Practice**: Use domain-style paths for public modules

---

## General Recovery Steps

### 1. Check Logs

Enable debug logging:

```bash
# Set log level to debug
export LOG_LEVEL=debug

# Run command and check output
```

### 2. Verify Prerequisites

```bash
# Check Go version
go version  # Should be 1.19+

# Check goctl
goctl --version

# Check environment
echo $GOPATH
echo $GOCTL_PATH
```

### 3. Clean Build

```bash
cd yourservice

# Clean build artifacts
go clean -cache -modcache -i -r

# Rebuild
go mod tidy
go build
```

### 4. Reset and Retry

```bash
# Remove generated service
rm -rf yourservice

# Try generation again with fresh state
```

### 5. Manual Fixes

If automatic fixes fail:

```bash
cd yourservice

# Fix imports
gofmt -w .
goimports -w .

# Fix modules
go mod init <module-name>
go mod tidy

# Verify build
go build
```

---

## Getting Help

### Check Documentation

1. **README.md** - Basic usage and examples
2. **CONTRIBUTING.md** - Development guidelines
3. **Issue Tracker** - Known issues and solutions

### Enable Verbose Output

Ask Claude for more details:

```
Show me the full error output and logs
```

### Report Issues

If error persists:

1. Check if it's a known issue on GitHub
2. Create a new issue with:
   - Error message
   - Steps to reproduce
   - Environment details (OS, Go version, goctl version)
   - Relevant logs

### Community Support

- GitHub Discussions
- go-zero Slack channel
- Stack Overflow (tag: go-zero)

---

## Prevention Best Practices

### 1. Validate Before Running

```
# Check service name before creation
# mcp-zero validates automatically

# Verify port availability
# mcp-zero checks automatically
```

### 2. Use Sensible Defaults

```
# Let mcp-zero choose defaults
Create API service "userservice"  # Uses default port 8888

# Or specify custom values
Create API service "userservice" on port 9000
```

### 3. Keep Tools Updated

```bash
# Update goctl regularly
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Update mcp-zero
go install github.com/zeromicro/mcp-zero@latest
```

### 4. Test in Isolation

```bash
# Test in temporary directory first
mkdir -p /tmp/test-service
cd /tmp/test-service
# Create service and verify it works
```

### 5. Backup Before Modifications

```bash
# Backup existing service before adding features
cp -r yourservice yourservice.backup
```

---

## Quick Reference

| Error Type | Quick Fix | Prevention |
|-----------|-----------|------------|
| Invalid service name | Use alphanumeric + underscores | Validate names first |
| Port in use | Change port or stop process | Check ports before use |
| goctl not found | Install goctl, set GOCTL_PATH | Install prerequisites |
| Build fails | Run go mod tidy | Let mcp-zero handle builds |
| Import errors | Run goimports | Let mcp-zero fix imports |
| DB connection | Check connection string | Test connection first |
| Config error | Fix YAML syntax | Use generated templates |

---

## Conclusion

Most errors have automatic fixes built into mcp-zero. When manual intervention is needed, this guide provides step-by-step solutions. Always start with automatic recovery, then progress to manual fixes if needed.

**Remember**: mcp-zero handles most common issues automatically through:
- Input validation
- Import fixing
- Module initialization
- Build verification
- Error recovery

For persistent issues, consult the community or file a bug report.
