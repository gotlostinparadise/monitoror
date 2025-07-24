//go:build faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

type DNSParams struct {
	params.Default

	RecordType      string `json:"recordType" query:"recordType"`
	Name            string `json:"name" query:"name"`
	ExpectedValue   string `json:"expectedValue,omitempty" query:"expectedValue"`
	ExpectedPattern string `json:"expectedPattern,omitempty" query:"expectedPattern"`

	Status coreModels.TileStatus `json:"status" query:"status"`
}

func (p *DNSParams) Validate() (errors []validator.Error) {
	if p.ExpectedPattern != "" {
		if _, err := regexp.Compile(p.ExpectedPattern); err != nil {
			errors = append(errors, validator.NewDefaultError("expectedPattern", "valid golang regex"))
		}
	}
	return
}
