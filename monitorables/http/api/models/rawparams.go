//go:build !faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/internal/pkg/validator"
)

type (
	HTTPRawParams struct {
		URL           string            `json:"url" query:"url" validate:"required,url,http"`
		Regex         string            `json:"regex,omitempty" query:"regex" validate:"regex"`
		StatusCodeMin *int              `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int              `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
		SSLVerify     *bool             `json:"sslVerify,omitempty" query:"sslVerify"`
		Headers       map[string]string `json:"headers,omitempty" query:"headers"`
	}
)

func (p *HTTPRawParams) Validate() []validator.Error {
	return validateStatusCode(p)
}

func (p *HTTPRawParams) GetURL() (url string) { return p.URL }
func (p *HTTPRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPRawParams) GetRegex() string              { return p.Regex }
func (p *HTTPRawParams) GetRegexp() *regexp.Regexp     { return getRegexp(p.GetRegex()) }
func (p *HTTPRawParams) GetSSLVerify() *bool           { return p.SSLVerify }
func (p *HTTPRawParams) GetHeaders() map[string]string { return p.Headers }
