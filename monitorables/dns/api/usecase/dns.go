//go:build !faker

package usecase

import (
	"fmt"
	"regexp"
	"strings"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/dns/api"
	"github.com/monitoror/monitoror/monitorables/dns/api/models"
)

type dnsUsecase struct {
	repository api.Repository
}

func NewDNSUsecase(repository api.Repository) api.Usecase {
	return &dnsUsecase{repository}
}

func (du *dnsUsecase) DNS(params *models.DNSParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.DNSTileType)
	tile.Label = fmt.Sprintf("%s %s", params.RecordType, params.Name)

	records, err := du.repository.Lookup(params.RecordType, params.Name)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		return tile, nil
	}

	match := false
	if params.ExpectedValue != "" {
		for _, r := range records {
			if r == params.ExpectedValue {
				match = true
				break
			}
		}
	}
	if !match && params.ExpectedPattern != "" {
		re, _ := regexp.Compile(params.ExpectedPattern)
		for _, r := range records {
			if re.MatchString(r) {
				match = true
				break
			}
		}
	}

	msg := ""
	if match {
		tile.Status = coreModels.SuccessStatus
		msg = fmt.Sprintf("%s record is %s as expected", params.RecordType, strings.Join(records, ","))
	} else {
		tile.Status = coreModels.FailedStatus
		msg = fmt.Sprintf("%s record changed to %s", params.RecordType, strings.Join(records, ","))
	}

	tile.WithMetrics(coreModels.RawUnit)
	tile.Metrics.Values = []string{msg}

	return tile, nil
}
