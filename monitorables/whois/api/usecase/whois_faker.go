//go:build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api"
	"github.com/monitoror/monitoror/monitorables/whois/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type whoisUsecase struct {
	refByDomain map[string]time.Time
}

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.WarningStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

func NewWHOISUsecase() api.Usecase { return &whoisUsecase{make(map[string]time.Time)} }

func (wu *whoisUsecase) WHOIS(params *models.WHOISParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.WHOISTileType).WithMetrics(coreModels.RawUnit)
	tile.Label = params.Domain

	status := wu.computeStatus(params.Domain)
	tile.Status = nonempty.Struct(params.Status, status).(coreModels.TileStatus)

	remaining := wu.computeRemainingDays(params.Domain)
	tile.Metrics.Values = []string{fmt.Sprintf("Expires in %d days", remaining)}
	tile.Message = time.Now().AddDate(0, 0, remaining).Format(time.RFC3339)

	return tile, nil
}

func (wu *whoisUsecase) computeStatus(domain string) coreModels.TileStatus {
	value, ok := wu.refByDomain[domain]
	if !ok {
		wu.refByDomain[domain] = faker.GetRefTime()
		value = wu.refByDomain[domain]
	}
	return faker.ComputeStatus(value, availableStatuses)
}

func (wu *whoisUsecase) computeRemainingDays(domain string) int {
	value, ok := wu.refByDomain[domain]
	if !ok {
		wu.refByDomain[domain] = faker.GetRefTime()
		value = wu.refByDomain[domain]
	}
	duration := faker.ComputeDuration(value, time.Hour*24*90)
	remaining := 90 - int(duration.Hours()/24)
	if remaining == 0 {
		remaining = 90
	}
	return remaining
}
