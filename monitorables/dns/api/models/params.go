//go:build !faker

package models

import (
    "regexp"

    "github.com/monitoror/monitoror/internal/pkg/monitorable/params"
    "github.com/monitoror/monitoror/internal/pkg/validator"
)

type DNSParams struct {
    params.Default

    RecordType      string `json:"recordType" query:"recordType" validate:"required,oneof=A AAAA CNAME TXT"`
    Name            string `json:"name" query:"name" validate:"required"`
    ExpectedValue   string `json:"expectedValue,omitempty" query:"expectedValue"`
    ExpectedPattern string `json:"expectedPattern,omitempty" query:"expectedPattern"`
}

func (p *DNSParams) Validate() (errors []validator.Error) {
    if p.ExpectedValue == "" && p.ExpectedPattern == "" {
        errors = append(errors, validator.Error{Field: "expectedValue", Tag: "required_without_expectedPattern"})
    }
    if p.ExpectedPattern != "" {
        if _, err := regexp.Compile(p.ExpectedPattern); err != nil {
            errors = append(errors, validator.Error{Field: "expectedPattern", Tag: "regexp"})
        }
    }
    return
}
