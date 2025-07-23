//go:generate mockery --name Repository

package api

import "time"

type (
	Certificate struct {
		NotBefore time.Time
		NotAfter  time.Time
		Issuer    string
		Subject   string
	}

	Repository interface {
		FetchCertificate(hostname string, port int) (*Certificate, error)
	}
)
