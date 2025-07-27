package usecase

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func buildMessage(display, domain string, expiry time.Time, raw string) string {
	d := strings.TrimSpace(strings.ToLower(display))

	if d == "" || d == "none" {
		return ""
	}

	if d == "full" {
		return fmt.Sprintf("Domain: %s / Expiration: %s", domain, expiry.Format(time.RFC3339))
	}

	parts := []string{}
	for _, field := range strings.Split(display, ",") {
		f := strings.TrimSpace(strings.ToLower(field))
		switch f {
		case "domain":
			parts = append(parts, domain)
		case "expiration", "expiry", "expires":
			parts = append(parts, expiry.Format(time.RFC3339))
		default:
			if f == "" {
				continue
			}
			if re, err := regexp.Compile(field); err == nil {
				if m := re.FindStringSubmatch(raw); m != nil {
					if len(m) > 1 {
						parts = append(parts, m[1])
					} else {
						parts = append(parts, m[0])
					}
				} else {
					parts = append(parts, "поле не найдено")
				}
			}
		}
	}

	return strings.Join(parts, " / ")
}
