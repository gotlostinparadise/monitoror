package http

import (
    "net/http"

    "github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
    "github.com/monitoror/monitoror/monitorables/dns/api"
    "github.com/monitoror/monitoror/monitorables/dns/api/models"

    "github.com/labstack/echo/v4"
)

type DNSDelivery struct {
    usecase api.Usecase
}

func NewDNSDelivery(u api.Usecase) *DNSDelivery { return &DNSDelivery{u} }

func (d *DNSDelivery) GetDNS(c echo.Context) error {
    params := &models.DNSParams{}
    if err := delivery.BindAndValidateParams(c, params); err != nil {
        return err
    }

    tile, err := d.usecase.DNS(params)
    if err != nil {
        return err
    }

    return c.JSON(http.StatusOK, tile)
}
