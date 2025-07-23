//go:build !faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/internal/pkg/validator"
)

type (
	HTTPFormattedParams struct {
		URL           string            `json:"url" query:"url" validate:"required,url,http"`
		Format        Format            `json:"format" query:"format" validate:"required,oneof=JSON YAML XML"`
		Key           string            `json:"key" query:"key" validate:"required,ne=."`
		Regex         string            `json:"regex,omitempty" query:"regex" validate:"regex"`
		StatusCodeMin *int              `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int              `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
		SSLVerify     *bool             `json:"sslVerify,omitempty" query:"sslVerify"`
		Headers       map[string]string `json:"headers,omitempty" query:"headers"`
	}
)

func (p *HTTPFormattedParams) Validate() []validator.Error {
	return validateStatusCode(p)
}

func (p *HTTPFormattedParams) GetURL() (url string) { return p.URL }
func (p *HTTPFormattedParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPFormattedParams) GetRegex() string          { return p.Regex }
func (p *HTTPFormattedParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPFormattedParams) GetSSLVerify() *bool           { return p.SSLVerify }
func (p *HTTPFormattedParams) GetHeaders() map[string]string { return p.Headers }

func (p *HTTPFormattedParams) GetKey() string    { return p.Key }
func (p *HTTPFormattedParams) GetFormat() Format { return p.Format }
