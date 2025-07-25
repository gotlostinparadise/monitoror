package http

import (
	"net/http"

	delivery "github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	"github.com/monitoror/monitoror/monitorables/rtsp/api/models"

	echo "github.com/labstack/echo/v4"
)

type RTSPDelivery struct {
	usecase api.Usecase
}

func NewRTSPDelivery(u api.Usecase) *RTSPDelivery { return &RTSPDelivery{u} }

func (h *RTSPDelivery) GetRTSP(c echo.Context) error {
	params := &models.RTSPParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.usecase.RTSP(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
