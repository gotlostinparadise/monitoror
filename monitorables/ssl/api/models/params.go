//go:build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	SSLParams struct {
		params.Default

		Hostname string `json:"hostname" query:"hostname" validate:"required"`
		Port     int    `json:"port" query:"port" validate:"required,gt=0"`
		WarnDays int    `json:"warnDays" query:"warnDays" validate:"required,gte=0"`
	}
)
