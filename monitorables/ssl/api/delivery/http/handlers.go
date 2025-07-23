package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"

	"github.com/labstack/echo/v4"
)

type SSLDelivery struct {
	sslUsecase api.Usecase
}

func NewSSLDelivery(u api.Usecase) *SSLDelivery {
	return &SSLDelivery{u}
}

func (h *SSLDelivery) GetSSL(c echo.Context) error {
	params := &models.SSLParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.sslUsecase.SSL(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
