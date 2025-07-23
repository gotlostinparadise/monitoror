package repository

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/config"
)

type sslRepository struct {
	config *config.SSL
}

func NewSSLRepository(conf *config.SSL) api.Repository {
	return &sslRepository{conf}
}

func (r *sslRepository) FetchCertificate(hostname string, port int) (*api.Certificate, error) {
	addr := fmt.Sprintf("%s:%d", hostname, port)
	dialer := &net.Dialer{Timeout: time.Millisecond * time.Duration(r.config.Timeout)}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no certificate")
	}

	cert := state.PeerCertificates[0]
	return &api.Certificate{
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
		Issuer:    cert.Issuer.String(),
		Subject:   cert.Subject.String(),
	}, nil
}
