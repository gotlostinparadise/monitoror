//go:build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	"github.com/monitoror/monitoror/monitorables/rtsp/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type rtspUsecase struct {
	timeRefByHost map[string]time.Time
}

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

func NewRTSPUsecase(_ api.Repository) api.Usecase { return &rtspUsecase{make(map[string]time.Time)} }

func (u *rtspUsecase) RTSP(params *models.RTSPParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.RTSPTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	tile.Status = nonempty.Struct(params.Status, u.computeStatus(params)).(coreModels.TileStatus)
	return tile, nil
}

func (u *rtspUsecase) computeStatus(params *models.RTSPParams) coreModels.TileStatus {
	key := fmt.Sprintf("%s:%d", params.Hostname, params.Port)
	value, ok := u.timeRefByHost[key]
	if !ok {
		u.timeRefByHost[key] = faker.GetRefTime()
		value = u.timeRefByHost[key]
	}
	return faker.ComputeStatus(value, availableStatuses)
}
