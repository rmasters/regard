package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"regard/internal/domain"
	"regard/internal/query"
)

// OutputJSON renders a QueryResult as formatted JSON
func OutputJSON(result query.QueryResult, useColor bool) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		return
	}

	if !useColor {
		fmt.Print(string(jsonBytes))
		return
	}

	highlightJSON(string(jsonBytes))
}

// OutputSummaryJSON renders a domain summary as formatted JSON
func OutputSummaryJSON(summary domain.Summary, useColor bool) {
	jsonBytes, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		return
	}

	if !useColor {
		fmt.Print(string(jsonBytes))
		return
	}

	highlightJSON(string(jsonBytes))
}

func highlightJSON(jsonStr string) {
	// Apply syntax highlighting
	lexer := lexers.Get("json")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, jsonStr)
	if err != nil {
		fmt.Print(jsonStr)
		return
	}

	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		fmt.Print(jsonStr)
		return
	}
}