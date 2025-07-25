//go:build !faker

package models

import "github.com/monitoror/monitoror/internal/pkg/monitorable/params"

type RTSPParams struct {
	params.Default

	Hostname string `json:"hostname" query:"hostname" validate:"required"`
	Port     int    `json:"port" query:"port" validate:"required,gt=0"`
	Path     string `json:"path,omitempty" query:"path"`
	Username string `json:"username" query:"username" validate:"required"`
	Password string `json:"password" query:"password" validate:"required"`
	Method   string `json:"method,omitempty" query:"method"`
}

func (p *RTSPParams) GetMethod() string {
	if p.Method == "" {
		return "OPTIONS"
	}
	return p.Method
}
