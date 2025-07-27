package models

import "github.com/monitoror/monitoror/internal/pkg/monitorable/params"

type CommandParams struct {
	params.Default

	Command     string `json:"command" query:"command" validate:"required"`
	ExitCodeMax *int   `json:"exitCodeMax,omitempty" query:"exitCodeMax"`
}

func (p *CommandParams) GetExitCodeMax() int {
	if p.ExitCodeMax != nil {
		return *p.ExitCodeMax
	}
	return 0
}
