//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type RTSPParams struct {
	params.Default

	Hostname string `json:"hostname" query:"hostname"`
	Port     int    `json:"port" query:"port"`
	Path     string `json:"path,omitempty" query:"path"`
	Username string `json:"username" query:"username"`
	Password string `json:"password" query:"password"`
	Method   string `json:"method,omitempty" query:"method"`

	Status coreModels.TileStatus `json:"status" query:"status"`
}

func (p *RTSPParams) GetMethod() string {
	if p.Method == "" {
		return "OPTIONS"
	}
	return p.Method
}
