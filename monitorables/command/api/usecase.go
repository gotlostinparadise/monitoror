//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api/models"
)

const (
	CommandTileType coreModels.TileType = "COMMAND-STATUS"
)

type Usecase interface {
	CommandStatus(params *models.CommandParams) (*coreModels.Tile, error)
}
