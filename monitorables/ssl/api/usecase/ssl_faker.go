//go:build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type sslUsecase struct {
	timeRefByHost map[string]time.Time
}

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.WarningStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

func NewSSLUsecase() api.Usecase {
	return &sslUsecase{make(map[string]time.Time)}
}

func (su *sslUsecase) SSL(params *models.SSLParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.SSLTileType).WithMetrics(coreModels.RawUnit)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	status := su.computeStatus(params)
	tile.Status = nonempty.Struct(params.Status, status).(coreModels.TileStatus)
	tile.Metrics.Values = []string{"", "", "issuer", "subject"}
	return tile, nil
}

func (su *sslUsecase) computeStatus(params *models.SSLParams) coreModels.TileStatus {
	value, ok := su.timeRefByHost[params.Hostname]
	if !ok {
		su.timeRefByHost[params.Hostname] = faker.GetRefTime()
	}
	return faker.ComputeStatus(value, availableStatuses)
}
