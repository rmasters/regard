# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**regard** is a command-line tool for domain research and discovery. It supports both traditional WHOIS protocol and the modern Registration Data Access Protocol (RDAP), making it ideal for domain hunters, researchers, and anyone investigating domain availability and ownership.

## Current State

The project has been fully implemented with proper Go project structure:
- **Well-structured codebase** following Go conventions
- **Modular architecture** with separation of concerns
- Support for domains, IP addresses, and ASN lookups
- RDAP-first with WHOIS fallback strategy
- Multiple output formats: human-readable summary, JSON, raw, and verbose
- Color-coded terminal output with syntax highlighting
- Built-in binary (`./regard`) ready for use

## Technology Stack

This is a **Go** application using Go 1.24.5 with the following key dependencies:

### Dependencies (from go.mod)
- `github.com/likexian/whois` v1.15.6 - WHOIS protocol client
- `github.com/openrdap/rdap` v0.9.1 - RDAP protocol client  
- `github.com/alecthomas/chroma/v2` v2.20.0 - Terminal syntax highlighting and coloring
- `golang.org/x/term` v0.34.0 - Terminal width detection

## Project Structure

```
regard/
├── cmd/
│   └── regard/
│       └── main.go              # CLI entry point and flag handling
├── internal/
│   ├── query/
│   │   ├── types.go            # QueryResult and types
│   │   ├── detector.go         # Query type detection
│   │   ├── rdap.go            # RDAP protocol implementation
│   │   └── whois.go           # WHOIS protocol implementation
│   ├── domain/
│   │   ├── types.go           # Domain summary types
│   │   ├── status.go          # Status interpretation logic
│   │   ├── guidance.go        # Post-expiration guidance
│   │   └── parser.go          # RDAP/WHOIS data parsing
│   └── output/
│       ├── json.go           # JSON formatting with syntax highlighting
│       ├── terminal.go       # Human-readable terminal output
│       └── usage.go          # Help text and usage information
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md
```

## Architecture

The refactored application follows proper Go project organization:

### Core Features
- **Multi-protocol support**: RDAP (primary) with WHOIS fallback
- **Query type detection**: Automatic detection of domain, IP (v4/v6), and ASN queries
- **Multiple output modes**:
  - Human-readable summary (default)
  - JSON summary (`--json`)
  - Full verbose output (`--v`)
  - Raw protocol output (`--raw`)
- **Color-coded output** with `--no-color` option
- **Protocol forcing**: `--rdap` or `--whois` flags
- **Terminal-aware formatting** with width detection

### Package Responsibilities
- **`cmd/regard`**: CLI interface, flag parsing, and main application flow
- **`internal/query`**: Protocol implementations, query execution, type detection
- **`internal/domain`**: Domain data modeling, status interpretation, guidance logic
- **`internal/output`**: All output formatting (terminal, JSON, usage)

### Key Benefits of Structure
- **Testable**: Each package can be independently unit tested
- **Maintainable**: Clear separation of concerns makes changes easier
- **Reusable**: Internal packages can be composed in different ways
- **Standard**: Follows Go community conventions with `cmd/` and `internal/`

### Domain Hunter Focus
- **Expiration guidance**: Oriented toward domain acquisition rather than ownership management
- **Availability detection**: Smart detection of domains available for registration
- **Drop monitoring**: Information about domain drop timelines and grace periods
- **Research-friendly**: Clean JSON output perfect for domain research tools and scripts

## Development Commands

```bash
# Build the application
go build -o regard cmd/regard/main.go

# Run directly during development
go run cmd/regard/main.go <domain>

# Test with different options
./regard example.com               # Default human-readable
./regard --json example.com        # JSON summary  
./regard -v example.com           # Full verbose output
./regard --raw example.com        # Raw protocol response
./regard --whois example.com      # Force WHOIS
./regard --rdap example.com       # Force RDAP only
./regard 8.8.8.8                 # IP lookup
./regard AS15169                  # ASN lookup
```

## Output Design

The tool follows httpie's philosophy with:
- **Clean terminal output** with color coding
- **Status-first display**: domain status prominently shown
- **Structured information**: timeline, nameservers, DNSSEC, registrar
- **Smart guidance**: expiration warnings and renewal advice
- **Flexible formats**: human-readable, JSON, or raw protocol data
- **Terminal integration**: width-aware formatting and color detection