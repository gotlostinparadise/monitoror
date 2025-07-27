package usecase

import (
	"fmt"
	"strings"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	"github.com/monitoror/monitoror/monitorables/command/api/models"
)

type commandUsecase struct {
	repository api.Repository
}

func NewCommandUsecase(r api.Repository) api.Usecase { return &commandUsecase{r} }

func (cu *commandUsecase) CommandStatus(params *models.CommandParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.CommandTileType)
	tile.Label = params.Command

	output, exitCode, _, err := cu.repository.Exec(params.Command)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		tile.Message = err.Error()
		return tile, nil
	}

	if exitCode <= params.GetExitCodeMax() {
		tile.Status = coreModels.SuccessStatus
	} else {
		tile.Status = coreModels.FailedStatus
	}
	tile.Message = strings.TrimSpace(output)
	return tile, nil
}
