//go:build !faker

package usecase

import (
	"encoding/hex"
	"fmt"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/models"
)

type (
	portUsecase struct {
		repository api.Repository
	}
)

func NewPortUsecase(repository api.Repository) api.Usecase {
	return &portUsecase{repository}
}

func (pu *portUsecase) Port(params *models.PortParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	var payload []byte
	if params.Payload != "" {
		payload, err = hex.DecodeString(params.Payload)
		if err != nil {
			tile.Status = coreModels.FailedStatus
			err = nil
			return
		}
	}

	err = pu.repository.OpenSocket(params.Hostname, params.Port, string(params.GetType()), payload)
	if err == nil {
		tile.Status = coreModels.SuccessStatus
	} else {
		tile.Status = coreModels.FailedStatus
		err = nil
	}

	return
}
