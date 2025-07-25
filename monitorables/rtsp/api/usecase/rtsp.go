//go:build !faker

package usecase

import (
	"fmt"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	"github.com/monitoror/monitoror/monitorables/rtsp/api/models"
)

type rtspUsecase struct {
	repository api.Repository
}

func NewRTSPUsecase(repo api.Repository) api.Usecase { return &rtspUsecase{repo} }

func (u *rtspUsecase) RTSP(params *models.RTSPParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.RTSPTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	success, duration, err := u.repository.Authenticate(params.Hostname, params.Port, params.Path, params.GetMethod(), params.Username, params.Password)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		return tile, nil
	}

	if success {
		tile.Status = coreModels.SuccessStatus
		tile.Message = "authenticated"
	} else {
		tile.Status = coreModels.FailedStatus
		tile.Message = "authentication failed"
	}

	tile.WithMetrics(coreModels.MillisecondUnit)
	tile.Metrics.Values = []string{fmt.Sprintf("%d", duration.Milliseconds())}

	return tile, nil
}
