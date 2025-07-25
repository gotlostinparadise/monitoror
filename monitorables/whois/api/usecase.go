//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api/models"
)

const (
	WHOISTileType coreModels.TileType = "WHOIS"
)

type Usecase interface {
	WHOIS(params *models.WHOISParams) (*coreModels.Tile, error)
}
