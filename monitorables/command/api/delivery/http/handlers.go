package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/api/models"

	"github.com/labstack/echo/v4"
)

type CommandDelivery struct {
	usecase api.Usecase
}

func NewCommandDelivery(u api.Usecase) *CommandDelivery { return &CommandDelivery{u} }

func (h *CommandDelivery) GetCommandStatus(c echo.Context) error {
	params := &models.CommandParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.usecase.CommandStatus(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
