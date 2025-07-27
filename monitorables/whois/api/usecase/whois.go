//go:build !faker

package usecase

import (
	"fmt"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api"
	"github.com/monitoror/monitoror/monitorables/whois/api/models"
)

type whoisUsecase struct {
	repository api.Repository
}

func NewWHOISUsecase(repository api.Repository) api.Usecase {
	return &whoisUsecase{repository}
}

func (wu *whoisUsecase) WHOIS(params *models.WHOISParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.WHOISTileType)
	tile.Label = params.Domain

	expiry, raw, err := wu.repository.DomainExpiration(params.Domain)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		tile.Message = err.Error()
		return tile, nil
	}

	remaining := int(expiry.Sub(time.Now()).Hours() / 24)
	if params.WarnDays > 0 && remaining <= params.WarnDays {
		tile.Status = coreModels.WarningStatus
	} else {
		tile.Status = coreModels.SuccessStatus
	}

	tile.WithMetrics(coreModels.RawUnit)
	tile.Metrics.Values = []string{fmt.Sprintf("Expires in %d days", remaining)}
	tile.Message = buildMessage(params.Display, params.Domain, expiry, raw)

	return tile, nil
}
