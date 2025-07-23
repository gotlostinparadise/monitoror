package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestSSLParams_Validate(t *testing.T) {
	param := &SSLParams{}
	test.AssertParams(t, param, 3)

	param = &SSLParams{Hostname: "test"}
	test.AssertParams(t, param, 2)

	param = &SSLParams{Hostname: "test", Port: 443}
	test.AssertParams(t, param, 1)

	param = &SSLParams{Hostname: "test", Port: 443, WarnDays: 10}
	test.AssertParams(t, param, 0)
}
