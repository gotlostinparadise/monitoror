//go:generate mockery --name Repository

package api

import "time"

type (
	Repository interface {
		Authenticate(hostname string, port int, path, method, username, password string) (bool, time.Duration, error)
	}
)
