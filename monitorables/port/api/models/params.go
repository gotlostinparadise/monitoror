//go:build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	PortType string

	PortParams struct {
		params.Default

		Hostname string   `json:"hostname" query:"hostname" validate:"required"`
		Port     int      `json:"port" query:"port" validate:"required,gt=0"`
		Type     PortType `json:"type,omitempty" query:"type" validate:"omitempty,oneof=tcp udp"`
		Payload  string   `json:"payload,omitempty" query:"payload" validate:"omitempty,hexadecimal"`
	}
)

const (
	TCPPortType PortType = "tcp"
	UDPPortType PortType = "udp"
)

func (p *PortParams) GetType() PortType {
	if p.Type == "" {
		return TCPPortType
	}
	return p.Type
}
