package output

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/term"

	"regard/internal/domain"
)

// OutputSummary renders a domain summary in human-readable format
func OutputSummary(summary domain.Summary, useColor bool) {
	// Color functions
	bold := func(s string) string {
		if useColor {
			return fmt.Sprintf("\033[1m%s\033[0m", s)
		}
		return s
	}

	green := func(s string) string {
		if useColor {
			return fmt.Sprintf("\033[32m%s\033[0m", s)
		}
		return s
	}

	yellow := func(s string) string {
		if useColor {
			return fmt.Sprintf("\033[33m%s\033[0m", s)
		}
		return s
	}

	red := func(s string) string {
		if useColor {
			return fmt.Sprintf("\033[31m%s\033[0m", s)
		}
		return s
	}

	blue := func(s string) string {
		if useColor {
			return fmt.Sprintf("\033[34m%s\033[0m", s)
		}
		return s
	}

	// Header - compact format: domain (status) <spacer> protocol
	statusColor := green
	statusText := summary.Status

	if summary.Status == "expired" {
		statusColor = red
	} else if summary.Status == "unknown" {
		statusColor = yellow
	} else if summary.Status == "available" {
		statusColor = func(s string) string {
			if useColor {
				return fmt.Sprintf("\033[1;32m%s\033[0m", strings.ToUpper(s)) // Bold green uppercase
			}
			return strings.ToUpper(s)
		}
		statusText = "AVAILABLE"
	}

	// Calculate padding for alignment using terminal width
	headerLeft := fmt.Sprintf("%s %s", bold(summary.Domain), statusColor(statusText))
	// Strip ANSI codes for length calculation
	headerLeftStripped := stripAnsiCodes(fmt.Sprintf("%s %s", summary.Domain, statusText))

	// Get terminal width, fallback to 80 if unable to detect
	termWidth := getTerminalWidth()
	rightSide := summary.Protocol

	// Calculate padding: total width - left side - right side - 1 space minimum
	padding := termWidth - len(headerLeftStripped) - len(rightSide) - 1
	if padding < 1 {
		padding = 1
	}

	fmt.Printf("%s%s%s\n", headerLeft, strings.Repeat(" ", padding), rightSide)

	// For available domains, show a celebratory message and skip most sections
	if summary.Status == "available" {
		fmt.Printf("\n🎉 %s\n", green("This domain appears to be available for registration!"))
		return
	}

	// Timeline (only show header if there are timeline entries)
	hasTimelineEntries := summary.Timeline.Registration != nil || summary.Timeline.LastUpdated != nil || summary.Timeline.Expiration != nil
	if hasTimelineEntries {
		fmt.Printf("\n%s\n", bold("Timeline:"))
		if summary.Timeline.Registration != nil {
			fmt.Printf("  • %s: %s (%s)\n",
				bold("Registered"),
				summary.Timeline.Registration.Date.Format("2006-01-02"),
				blue(summary.Timeline.Registration.HumanReadable))
		}
		if summary.Timeline.LastUpdated != nil {
			fmt.Printf("  • %s: %s (%s)\n",
				bold("Last updated"),
				summary.Timeline.LastUpdated.Date.Format("2006-01-02"),
				blue(summary.Timeline.LastUpdated.HumanReadable))
		}
		if summary.Timeline.Expiration != nil {
			expiryColor := green
			if summary.Timeline.Expiration.Date.Before(time.Now()) {
				expiryColor = red
			} else if summary.Timeline.Expiration.Date.Before(time.Now().AddDate(0, 0, 30)) {
				expiryColor = yellow
			}
			fmt.Printf("  • %s: %s (%s)\n",
				bold("Expires"),
				summary.Timeline.Expiration.Date.Format("2006-01-02"),
				expiryColor(summary.Timeline.Expiration.HumanReadable))
		}
	}

	// Nameservers
	if len(summary.Nameservers) > 0 {
		fmt.Printf("\n%s\n", bold("Nameservers:"))
		for _, ns := range summary.Nameservers {
			fmt.Printf("  • %s\n", ns)
		}
	}

	// DNSSEC (only show for domains, not ASNs or IPs)
	if summary.QueryType == "domain" {
		fmt.Printf("\n%s ", bold("DNSSEC:"))
		if summary.DNSSEC.Enabled {
			fmt.Printf("%s", green("enabled"))
			if summary.DNSSEC.Details != "" {
				fmt.Printf(" (%s)", summary.DNSSEC.Details)
			}
		} else {
			fmt.Printf("%s", red("disabled"))
		}
		fmt.Println()
	}

	// Registrar
	if summary.Registrar.Name != "" {
		fmt.Printf("\n%s %s", bold("Registrar:"), summary.Registrar.Name)
		if summary.Registrar.ID != "" {
			fmt.Printf(" (ID: %s)", summary.Registrar.ID)
		}
		fmt.Println()
	}

	// ASN Information
	if summary.ASN != nil {
		fmt.Printf("\n%s %s", bold("Organization:"), summary.ASN.Organization)
		if summary.ASN.Country != "" {
			fmt.Printf(" (%s)", summary.ASN.Country)
		}
		fmt.Println()

		if summary.ASN.Name != "" {
			fmt.Printf("%s %s\n", bold("AS Name:"), summary.ASN.Name)
		}

		if summary.ASN.AbuseContact != "" {
			fmt.Printf("%s %s\n", bold("Abuse Contact:"), summary.ASN.AbuseContact)
		}

		if len(summary.ASN.Peers) > 0 {
			fmt.Printf("\n%s\n", bold("Key Peers:"))
			// Show first 10 peers to avoid overwhelming output
			maxPeers := 10
			if len(summary.ASN.Peers) < maxPeers {
				maxPeers = len(summary.ASN.Peers)
			}
			for i := 0; i < maxPeers; i++ {
				fmt.Printf("  • %s\n", summary.ASN.Peers[i])
			}
			if len(summary.ASN.Peers) > 10 {
				fmt.Printf("  ... and %d more\n", len(summary.ASN.Peers)-10)
			}
		}
	}

	// Domain status details
	if len(summary.StatusDetails) > 0 {
		fmt.Printf("\n%s\n", bold("Status details:"))
		for _, status := range summary.StatusDetails {
			fmt.Printf("  • %s\n", status)
		}
	}

	// Post-expiration guidance
	if summary.PostExpiration != nil {
		fmt.Printf("\n%s\n", bold("Post-expiration guidance:"))
		fmt.Printf("  %s\n", summary.PostExpiration.GuidanceMessage)
	}
}

func stripAnsiCodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}

func getTerminalWidth() int {
	// Try to get terminal width
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return width
	}

	// Fallback to 80 columns if unable to detect
	return 80
}
