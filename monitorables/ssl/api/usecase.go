//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"
)

const (
	SSLTileType coreModels.TileType = "SSL"
)

type (
	Usecase interface {
		SSL(params *models.SSLParams) (*coreModels.Tile, error)
	}
)
