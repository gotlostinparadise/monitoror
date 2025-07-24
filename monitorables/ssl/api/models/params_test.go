package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	"github.com/stretchr/testify/assert"
)

func TestSSLParams_Validate(t *testing.T) {
	param := &SSLParams{}
	test.AssertParams(t, param, 2)

	param = &SSLParams{Domain: "test"}
	test.AssertParams(t, param, 1)

	param = &SSLParams{Domain: "test", WarnDays: 10}
	test.AssertParams(t, param, 0)
}

func TestSSLParams_GetPort(t *testing.T) {
	p := &SSLParams{Domain: "example.com", WarnDays: 1}
	assert.Equal(t, 443, p.GetPort())

	p = &SSLParams{Domain: "example.com", Port: 8443, WarnDays: 1}
	assert.Equal(t, 8443, p.GetPort())
}
