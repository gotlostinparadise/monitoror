//go:build faker

package usecase

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/api/models"
)

type commandUsecase struct{}

func NewCommandUsecase(_ api.Repository) api.Usecase { return &commandUsecase{} }

func (cu *commandUsecase) CommandStatus(params *models.CommandParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.CommandTileType)
	tile.Label = params.Command
	tile.Status = params.GetStatus()
	tile.Message = params.GetMessage()
	return tile, nil
}
