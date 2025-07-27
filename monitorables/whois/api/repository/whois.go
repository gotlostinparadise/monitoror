package repository

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/whois/api"
)

type whoisRepository struct{}

func NewWHOISRepository() api.Repository { return &whoisRepository{} }

func (r *whoisRepository) query(server, query string) (string, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(server, "43"), time.Second*5)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if _, err = fmt.Fprintf(conn, "%s\r\n", query); err != nil {
		return "", err
	}
	data, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *whoisRepository) DomainExpiration(domain string) (time.Time, string, error) {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return time.Time{}, "", fmt.Errorf("invalid domain")
	}
	tld := parts[len(parts)-1]

	res, err := r.query("whois.iana.org", tld)
	if err != nil {
		return time.Time{}, "", err
	}
	server := ""
	scanner := bufio.NewScanner(strings.NewReader(res))
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if strings.HasPrefix(line, "whois:") {
			server = strings.TrimSpace(scanner.Text()[6:])
			break
		}
	}
	if server == "" {
		return time.Time{}, "", fmt.Errorf("whois server not found")
	}

	out, err := r.query(server, domain)
	if err != nil {
		return time.Time{}, "", err
	}

	// Look for any expire-related line. Examples include:
	// "Expiry Date: 2026-06-07" or "Record expires on: Fri Oct 21 05:54:20 2033"
	re := regexp.MustCompile(`(?i)^\s*(?:record|domain|registrar|registry)?\s*(?:registration\s+)?(?:expir(?:ation|es?|y|e)|paid[- ]?till)\b\s*(?:date|on)?\s*[:-]?\s*(.+?)\s*$`)

	scanner = bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := re.FindStringSubmatch(line); m != nil {
			dateStr := strings.TrimSpace(m[1])
			layouts := []string{
				time.RFC3339,
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05Z",
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02 15:04:05-07",
				"2006-01-02",
				"Mon Jan 2 15:04:05 2006", // handle “Fri Oct 21 05:54:20 2033” and leading zero days
				"Mon Jan 2 2006",          // handle “Wed Sep 09 2026”
			}
			for _, l := range layouts {
				if t, err := time.Parse(l, dateStr); err == nil {
					return t.UTC(), out, nil // normalise to UTC
				}
			}
		}
	}
	return time.Time{}, out, fmt.Errorf("expiration date not found")
}
