//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	PortType string

	PortParams struct {
		params.Default

		Hostname string   `json:"hostname" query:"hostname"`
		Port     int      `json:"port" query:"port"`
		Type     PortType `json:"type,omitempty" query:"type"`
		Payload  string   `json:"payload,omitempty" query:"payload"`
		Display  string   `json:"display,omitempty" query:"display"`

		Status coreModels.TileStatus `json:"status" query:"status"`
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
