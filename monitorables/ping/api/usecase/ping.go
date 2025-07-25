//go:build !faker

package usecase

import (
	"fmt"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"
)

type (
	pingUsecase struct {
		repository api.Repository
	}
)

func NewPingUsecase(repository api.Repository) api.Usecase {
	return &pingUsecase{repository}
}

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.PingTileType)
	tile.Label = params.Hostname

	ping, err := pu.repository.ExecutePing(params.Hostname)
	if err == nil {
		tile.Status = coreModels.SuccessStatus
		tile.WithMetrics(coreModels.MillisecondUnit)
		if ping.Average < time.Millisecond {
			ms := float64(ping.Average.Nanoseconds()) / float64(time.Millisecond)
			tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%.2f", ms))
		} else {
			tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", ping.Average.Milliseconds()))
		}
	} else {
		tile.Status = coreModels.FailedStatus
		err = nil
	}

	return
}
