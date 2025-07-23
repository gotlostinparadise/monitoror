//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	HTTPStatusParams struct {
		URL           string            `json:"url" query:"url" validate:"required,url,http"`
		StatusCodeMin *int              `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int              `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
		SSLVerify     *bool             `json:"sslVerify,omitempty" query:"sslVerify"`
		Headers       map[string]string `json:"headers,omitempty" query:"headers"`

		Status  coreModels.TileStatus `json:"status" query:"status"`
		Message string                `json:"message" query:"message"`
	}
)

func (p *HTTPStatusParams) Validate() []validator.Error {
	return validateStatusCode(p)
}

func (p *HTTPStatusParams) GetURL() (url string) { return p.URL }
func (p *HTTPStatusParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPStatusParams) GetSSLVerify() *bool           { return p.SSLVerify }
func (p *HTTPStatusParams) GetHeaders() map[string]string { return p.Headers }

func (p *HTTPStatusParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *HTTPStatusParams) GetMessage() string                      { return p.Message }
func (p *HTTPStatusParams) GetValueValues() []string                { panic("unimplemented") }
func (p *HTTPStatusParams) GetValueUnit() coreModels.TileValuesUnit { panic("unimplemented") }
