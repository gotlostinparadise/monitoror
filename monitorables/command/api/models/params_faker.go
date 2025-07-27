//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type CommandParams struct {
	params.Default

	Command     string `json:"command" query:"command" validate:"required"`
	ExitCodeMax *int   `json:"exitCodeMax,omitempty" query:"exitCodeMax"`
	Display     string `json:"display,omitempty" query:"display"`
	Metrics     string `json:"metrics,omitempty" query:"metrics"`

	Status  coreModels.TileStatus `json:"status" query:"status"`
	Message string                `json:"message" query:"message"`
}

func (p *CommandParams) GetExitCodeMax() int {
	if p.ExitCodeMax != nil {
		return *p.ExitCodeMax
	}
	return 0
}

func (p *CommandParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *CommandParams) GetMessage() string                      { return p.Message }
func (p *CommandParams) GetValueValues() []string                { return nil }
func (p *CommandParams) GetValueUnit() coreModels.TileValuesUnit { return coreModels.RawUnit }
