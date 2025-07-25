//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type WHOISParams struct {
	params.Default

	Domain   string `json:"domain" query:"domain"`
	WarnDays int    `json:"warnDays" query:"warnDays"`

	Status coreModels.TileStatus `json:"status" query:"status"`
}
