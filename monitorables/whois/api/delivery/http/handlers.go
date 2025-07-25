package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/whois/api"
	"github.com/monitoror/monitoror/monitorables/whois/api/models"

	"github.com/labstack/echo/v4"
)

type WHOISDelivery struct {
	whoisUsecase api.Usecase
}

func NewWHOISDelivery(u api.Usecase) *WHOISDelivery { return &WHOISDelivery{u} }

func (h *WHOISDelivery) GetWHOIS(c echo.Context) error {
	params := &models.WHOISParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.whoisUsecase.WHOIS(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
