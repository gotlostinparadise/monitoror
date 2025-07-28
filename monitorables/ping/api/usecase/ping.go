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

func (pu *pingUsecase) Ping(params *models.PingParams) (*coreModels.Tile, error) {
    tile := coreModels.NewTile(api.PingTileType)
    tile.Label = params.Hostname
    tile.WithMetrics(coreModels.MillisecondUnit)

    ping, err := pu.repository.ExecutePing(params.Hostname)
    if err != nil {
        tile.Status = coreModels.FailedStatus
        return tile, nil
    }

    tile.Status = coreModels.SuccessStatus
    d := ping.Average

    switch {
    case d < time.Millisecond:
        // show two decimals for sub-ms
        ms := d.Seconds() * 1_000
        tile.Metrics.Values = append(
            tile.Metrics.Values,
            fmt.Sprintf("%.2f", ms),
        )

    default:
        // round to nearest whole ms
        rounded := d.Round(time.Millisecond)
        tile.Metrics.Values = append(
            tile.Metrics.Values,
            fmt.Sprintf("%d", rounded.Milliseconds()),
        )
    }

    return tile, nil
}

