package http

import (
	netHttp "net/http"
	"strings"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"

	"github.com/labstack/echo/v4"
)

type HTTPDelivery struct {
	httpUsecase api.Usecase
}

func NewHTTPDelivery(p api.Usecase) *HTTPDelivery {
	return &HTTPDelivery{p}
}

func parseHeaders(values map[string][]string) map[string]string {
	headers := make(map[string]string)
	for k, v := range values {
		if strings.HasPrefix(k, "headers[") && strings.HasSuffix(k, "]") {
			name := k[len("headers[") : len(k)-1]
			if len(v) > 0 {
				headers[name] = v[len(v)-1]
			}
		}
	}
	if len(headers) == 0 {
		return nil
	}
	return headers
}

func (h *HTTPDelivery) GetHTTPStatus(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPStatusParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}
	params.Headers = parseHeaders(c.QueryParams())

	tile, err := h.httpUsecase.HTTPStatus(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPRaw(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPRawParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}
	params.Headers = parseHeaders(c.QueryParams())

	tile, err := h.httpUsecase.HTTPRaw(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPFormatted(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPFormattedParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}
	params.Headers = parseHeaders(c.QueryParams())

	tile, err := h.httpUsecase.HTTPFormatted(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}
