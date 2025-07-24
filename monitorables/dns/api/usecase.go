package api

import (
    coreModels "github.com/monitoror/monitoror/models"
    "github.com/monitoror/monitoror/monitorables/dns/api/models"
)

const (
    DNSTileType coreModels.TileType = "DNS"
)

type Usecase interface {
    DNS(params *models.DNSParams) (*coreModels.Tile, error)
}
