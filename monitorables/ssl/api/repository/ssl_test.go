package repository

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorables/ssl/config"
	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchCertificate_Success(t *testing.T) {
	ts := httptest.NewTLSServer(nil)
	defer ts.Close()

	hostPort := strings.TrimPrefix(ts.URL, "https://")
	parts := strings.Split(hostPort, ":")
	host := parts[0]
	port := 0
	fmt.Sscanf(parts[1], "%d", &port)

	repo := NewSSLRepository(&config.SSL{Timeout: 2000})
	cert, err := repo.FetchCertificate(host, port)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}
