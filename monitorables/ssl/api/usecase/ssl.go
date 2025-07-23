//go:build !faker

package usecase

import (
	"fmt"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	"github.com/monitoror/monitoror/monitorables/ssl/api/models"
)

type sslUsecase struct {
	repository api.Repository
}

func NewSSLUsecase(repository api.Repository) api.Usecase {
	return &sslUsecase{repository}
}

func (su *sslUsecase) SSL(params *models.SSLParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.SSLTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	cert, err := su.repository.FetchCertificate(params.Hostname, params.Port)
	if err != nil {
		tile.Status = coreModels.FailedStatus
		return tile, nil
	}

	remaining := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	if remaining <= params.WarnDays {
		tile.Status = coreModels.WarningStatus
	} else {
		tile.Status = coreModels.SuccessStatus
	}

	tile.WithMetrics(coreModels.RawUnit)
	tile.Metrics.Values = []string{
		cert.NotBefore.Format(time.RFC3339),
		cert.NotAfter.Format(time.RFC3339),
		cert.Issuer,
		cert.Subject,
	}

	return tile, nil
}
