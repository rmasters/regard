package main

import (
	"flag"
	"fmt"
	"os"

	"regard/internal/domain"
	"regard/internal/output"
	"regard/internal/query"
)

func main() {
	var (
		useWhois   = flag.Bool("whois", false, "Force use of WHOIS protocol")
		useRdap    = flag.Bool("rdap", false, "Force use of RDAP protocol")
		rawOutput  = flag.Bool("raw", false, "Output raw response without formatting")
		verbose    = flag.Bool("v", false, "Verbose output (full details)")
		jsonOutput = flag.Bool("json", false, "Output in JSON format")
		noColor    = flag.Bool("no-color", false, "Disable syntax highlighting")
		showHelp   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *showHelp || len(os.Args) < 2 {
		output.PrintUsage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No domain specified\n")
		output.PrintUsage()
		os.Exit(1)
	}

	queryStr := args[0]

	var result query.QueryResult

	// Try RDAP first unless WHOIS is explicitly requested
	if !*useWhois {
		result = query.PerformRDAPQuery(queryStr)
		if !result.Success && !*useRdap {
			// Fall back to WHOIS if RDAP fails and not forced to use RDAP only
			result = query.PerformWhoisQuery(queryStr)
		}
	} else {
		result = query.PerformWhoisQuery(queryStr)
	}

	// Output the result
	if *rawOutput {
		if result.RawData != "" {
			fmt.Print(result.RawData)
		} else {
			fmt.Printf("Error: %s\n", result.Error)
		}
	} else if *verbose {
		// Full JSON output for verbose mode
		output.OutputJSON(result, !*noColor)
	} else if *jsonOutput {
		// Summary in JSON format
		if result.Success {
			summary := domain.CreateSummary(result)
			output.OutputSummaryJSON(summary, !*noColor)
		} else {
			fmt.Printf("{\"error\": \"%s\"}\n", result.Error)
		}
	} else {
		// Default: human-readable summary
		if result.Success {
			summary := domain.CreateSummary(result)
			output.OutputSummary(summary, !*noColor)
		} else {
			fmt.Printf("Error: %s\n", result.Error)
		}
	}
}
