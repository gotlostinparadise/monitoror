//go:generate mockery --name Repository

package api

import "time"

type Repository interface {
	DomainExpiration(domain string) (time.Time, string, error)
}
