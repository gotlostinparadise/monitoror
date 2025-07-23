package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
)

func TestPortParams_Validate(t *testing.T) {
	param := &PortParams{}
	test.AssertParams(t, param, 2)

	param = &PortParams{Hostname: "test"}
	test.AssertParams(t, param, 1)

	param = &PortParams{Hostname: "test", Port: 22}
	test.AssertParams(t, param, 0)

	param = &PortParams{Hostname: "test", Port: 22, Type: "invalid"}
	test.AssertParams(t, param, 1)
}
