# Regard - httpie for domains

`regard` helps make domain hunting from the command-line easier, with clean httpie-style readable output rather than lengthy raw WHOIS records. It also supports the modern RDAP protocol.

## Features

- 🚀 **Modern protocols**: RDAP-first with WHOIS fallback
- 🎯 **Smart detection**: Automatically detects domains, IPs (v4/v6), and ASNs  
- 🎨 **Beautiful output**: Syntax-highlighted JSON and human-readable summaries
- ⚡ **Fast**: Efficient Go implementation with minimal dependencies
- 🔧 **Flexible**: Multiple output formats and protocol options
- 📊 **Comprehensive**: Timeline, nameservers, DNSSEC, registrar info, and drop predictions
- 🎯 **Domain hunting**: Expiration monitoring and availability guidance for domain acquisition

## Installation

### From Source

```bash
git clone https://github.com/your-username/regard.git
cd regard
go build -o regard cmd/regard/main.go
```

### Using Go Install

```bash
go install github.com/rmasters/regard/cmd/regard@latest
```

## Usage

```bash
# Human-readable domain summary (default)
$ regard example.com

example.com active                                                         RDAP

Timeline:
  • Registered: 1995-08-14 (30 years ago)
  • Last updated: 2025-09-06 (today)
  • Expires: 2026-08-13 (in 11 months)

Nameservers:
  • A.IANA-SERVERS.NET
  • B.IANA-SERVERS.NET

DNSSEC: enabled (Delegation signed)

Registrar: RESERVED-Internet Assigned Numbers Authority (ID: 376)

Status details:
  • client delete prohibited
  • client transfer prohibited
  • client update prohibited
```

### JSON summaries

```bash
$ regard --json google.com
{
  "domain": "google.com",
  "status": "active",
  "protocol": "RDAP",
  "timeline": {
    "registration": {
      "date": "1997-09-15T04:00:00Z",
      "human_readable": "27 years ago"
    },
    "expiration": {
      "date": "2028-09-14T04:00:00Z", 
      "human_readable": "in 3 years"
    }
  },
  "nameservers": ["NS1.GOOGLE.COM", "NS2.GOOGLE.COM", "NS3.GOOGLE.COM", "NS4.GOOGLE.COM"],
  "dnssec": {"enabled": false},
  "registrar": {"name": "MarkMonitor Inc.", "id": "292"}
}
```

### Exotic uses

```bash
# Full verbose JSON output
regard -v github.com

# Raw protocol response
regard --raw stackoverflow.com

# IP address lookup
regard 8.8.8.8

# ASN lookup  
regard AS15169

# Force a specific protocol
regard --whois example.com
regard --rdap example.com
```

### Command Line Options

```
USAGE:
    regard [OPTIONS] <domain|ip|asn>

OPTIONS:
    --whois        Force use of WHOIS protocol
    --rdap         Force use of RDAP protocol only
    -v             Verbose output (full details)
    --json         Output summary in JSON format
    --raw          Output raw response without JSON formatting
    --no-color     Disable syntax highlighting
    --help         Show this help message
```

## Supported Query Types

| Type | Examples | Description |
|------|----------|-------------|
| **Domains** | `example.com`, `sub.example.org` | Any domain name |
| **IPv4** | `8.8.8.8`, `192.168.1.1` | IPv4 addresses |
| **IPv6** | `2001:4860:4860::8888` | IPv6 addresses |
| **ASN** | `AS15169`, `AS13335` | Autonomous System Numbers |

### Protocol Selection

`regard` prefers RDAP, but also works with WHOIS:

1. **RDAP first**: Modern, structured JSON responses
2. **WHOIS fallback**: Traditional protocol when RDAP unavailable
3. **Manual override**: Use `--rdap` or `--whois` to force a specific protocol

### Output Formats

- **Default**: Clean, human-readable summary with colors
- **JSON** (`--json`): Structured summary perfect for scripts
- **Verbose** (`-v`): Complete protocol response as formatted JSON
- **Raw** (`--raw`): Unprocessed protocol response

### Domain Drop Intelligence
- Identifies domains approaching expiration for monitoring
- Provides TLD-specific drop timelines and grace period information  
- Estimates when expired domains will become available for registration
- Distinguishes between renewal grace, redemption, and pending delete phases

### DNSSEC Information
- Shows DNSSEC delegation status
- Displays signing details when available

## Project Structure

```
regard/
├── cmd/regard/          # CLI application entry point
├── internal/
│   ├── query/          # Protocol implementations (RDAP, WHOIS)
│   ├── domain/         # Domain logic and data modeling
│   └── output/         # Output formatting (terminal, JSON)
├── go.mod
├── LICENSE
└── README.md
```

## Development

### Building

```bash
go build -o regard cmd/regard/main.go
```

### Running Tests

```bash
go test ./...
```

## Acknowledgments

- Built with [github.com/likexian/whois](https://github.com/likexian/whois) for WHOIS protocol support
- Uses [github.com/openrdap/rdap](https://github.com/openrdap/rdap) for modern RDAP queries  
- Syntax highlighting powered by [github.com/alecthomas/chroma](https://github.com/alecthomas/chroma)
- Completely vibe-coded with Claude Code (Sonnet 4)