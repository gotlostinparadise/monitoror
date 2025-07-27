//go:build !faker

package models

import "github.com/monitoror/monitoror/internal/pkg/monitorable/params"

type WHOISParams struct {
	params.Default

	Domain   string `json:"domain" query:"domain" validate:"required"`
	WarnDays int    `json:"warnDays" query:"warnDays" validate:"gte=0"`
	Display  string `json:"display" query:"display"`
}
