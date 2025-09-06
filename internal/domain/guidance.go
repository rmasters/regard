package domain

import (
	"fmt"
	"strings"
	"time"
)

// GeneratePostExpirationGuidance provides guidance for domain hunters interested in expired domains
func GeneratePostExpirationGuidance(summary Summary) *ExpirationInfo {
	if summary.Timeline.Expiration == nil {
		return nil
	}

	now := time.Now()
	expiryDate := summary.Timeline.Expiration.Date

	// Only provide guidance if domain is expired or expiring soon
	if expiryDate.After(now.AddDate(0, 0, 30)) {
		// Domain expires in more than 30 days, no guidance needed
		return nil
	}

	daysUntilExpiry := int(expiryDate.Sub(now).Hours() / 24)
	daysExpired := int(now.Sub(expiryDate).Hours() / 24)

	guidance := &ExpirationInfo{}

	if expiryDate.After(now) {
		// Domain not yet expired - guidance for domain hunters
		if daysUntilExpiry <= 7 {
			guidance.GuidanceMessage = fmt.Sprintf("Domain expires in %d days. Monitor closely - it may become available if not renewed.", daysUntilExpiry)
		} else {
			guidance.GuidanceMessage = fmt.Sprintf("Domain expires in %d days. Add to your watchlist if interested.", daysUntilExpiry)
		}
		return guidance
	}

	// Domain is expired - guidance for acquisition
	guidance.DaysExpired = daysExpired

	// Extract TLD for specific guidance
	tld := extractTLD(summary.Domain)

	switch strings.ToLower(tld) {
	case "com", "net", "org":
		// Standard gTLD rules for domain hunters
		if daysExpired <= 30 {
			guidance.GuidanceMessage = "Domain is in renewal grace period. Original owner can still renew. Not yet available for registration."
		} else if daysExpired <= 75 {
			guidance.GuidanceMessage = "Domain is in redemption grace period. Original owner can still recover it with fees. Not available for public registration yet."
		} else {
			pendingDeleteDays := 80 - daysExpired
			if pendingDeleteDays > 0 {
				guidance.GuidanceMessage = fmt.Sprintf("Domain is pending deletion! It will drop and become available for registration in approximately %d days.", pendingDeleteDays)
				estimatedAvailable := expiryDate.AddDate(0, 0, 80)
				guidance.AvailableDate = &estimatedAvailable
			} else {
				guidance.GuidanceMessage = "Domain has completed the deletion process and should be available for registration at any registrar."
			}
		}
	case "uk", "co.uk":
		if daysExpired <= 90 {
			guidance.GuidanceMessage = "Domain is in renewal grace period (.uk domains have 90-day grace period). Original owner can still renew."
		} else {
			guidance.GuidanceMessage = "Domain has passed the renewal grace period and should be available for public registration."
		}
	default:
		// Generic guidance for other TLDs - domain hunter focused
		if daysExpired <= 30 {
			guidance.GuidanceMessage = "Domain may still be in renewal grace period. Original owner might still renew."
		} else {
			guidance.GuidanceMessage = "Domain is likely available for registration. Check with registrars or domain drop services."
		}
	}

	return guidance
}

func extractTLD(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		// Handle cases like co.uk, com.au
		if len(parts) >= 3 && len(parts[len(parts)-2]) <= 3 && len(parts[len(parts)-1]) <= 3 {
			return strings.Join(parts[len(parts)-2:], ".")
		}
		return parts[len(parts)-1]
	}
	return domain
}