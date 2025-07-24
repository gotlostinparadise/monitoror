//go:build faker

package usecase

import (
    "fmt"
    "time"

    "github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
    "github.com/monitoror/monitoror/pkg/nonempty"
    coreModels "github.com/monitoror/monitoror/models"
    "github.com/monitoror/monitoror/monitorables/dns/api"
    "github.com/monitoror/monitoror/monitorables/dns/api/models"
)

type dnsUsecase struct {
    refByName map[string]time.Time
}

var availableStatuses = faker.Statuses{
    {coreModels.SuccessStatus, time.Second * 30},
    {coreModels.FailedStatus, time.Second * 30},
}

func NewDNSUsecase() api.Usecase { return &dnsUsecase{make(map[string]time.Time)} }

func (du *dnsUsecase) DNS(params *models.DNSParams) (*coreModels.Tile, error) {
    tile := coreModels.NewTile(api.DNSTileType)
    tile.Label = fmt.Sprintf("%s %s", params.RecordType, params.Name)

    tile.Status = nonempty.Struct(params.Status, du.computeStatus(params)).(coreModels.TileStatus)
    return tile, nil
}

func (du *dnsUsecase) computeStatus(params *models.DNSParams) coreModels.TileStatus {
    key := fmt.Sprintf("%s-%s", params.RecordType, params.Name)
    value, ok := du.refByName[key]
    if !ok {
        du.refByName[key] = faker.GetRefTime()
        value = du.refByName[key]
    }
    return faker.ComputeStatus(value, availableStatuses)
}
