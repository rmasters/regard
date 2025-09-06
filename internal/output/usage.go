package output

import "fmt"

// PrintUsage displays the command usage information
func PrintUsage() {
	fmt.Printf(`regard - domain research and discovery tool

USAGE:
    regard [OPTIONS] <domain|ip|asn>

EXAMPLES:
    regard example.com          # Human-readable domain summary
    regard -v example.com       # Full verbose JSON output
    regard --json example.com   # Summary in JSON format
    regard --whois example.com  # Force WHOIS query
    regard --rdap example.com   # Force RDAP query only
    regard 8.8.8.8              # Query IP address
    regard AS15169              # Query ASN
    regard --raw example.com    # Raw output without formatting

OPTIONS:
    --whois        Force use of WHOIS protocol
    --rdap         Force use of RDAP protocol only
    -v             Verbose output (full details)
    --json         Output summary in JSON format
    --raw          Output raw response without JSON formatting
    --no-color     Disable syntax highlighting
    --help         Show this help message

By default, regard shows a human-readable summary and attempts RDAP first with WHOIS fallback.
`)
}
