//go:generate mockery --name Repository

package api

import "time"

type Repository interface {
	Exec(command string) (output string, exitCode int, duration time.Duration, err error)
}
