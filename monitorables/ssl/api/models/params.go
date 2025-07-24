//go:build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	SSLParams struct {
		params.Default

		Domain   string `json:"domain" query:"domain" validate:"required"`
		Port     int    `json:"port,omitempty" query:"port" validate:"omitempty,gt=0"`
		WarnDays int    `json:"warnDays" query:"warnDays" validate:"required,gte=0"`
		Display  string `json:"display,omitempty" query:"display"`
	}
)

func (p *SSLParams) GetPort() int {
	if p.Port == 0 {
		return 443
	}
	return p.Port
}
