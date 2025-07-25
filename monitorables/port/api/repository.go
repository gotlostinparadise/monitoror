//go:generate mockery --name Repository

package api

import "time"

type (
	Repository interface {
		// OpenSocket tries to connect and optionally send a payload.
		// It returns if the remote answered at least one byte, the received banner and the connection time.
		OpenSocket(hostname string, port int, network string, payload []byte) (bool, string, time.Duration, error)
	}
)
