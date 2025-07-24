//go:build !faker

package usecase

import (
	"fmt"
	"strings"
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
	tile.Label = fmt.Sprintf("%s:%d", params.Domain, params.GetPort())

	cert, err := su.repository.FetchCertificate(params.Domain, params.GetPort())
	if err != nil {
		tile.Status = coreModels.FailedStatus
		tile.Message = err.Error()
		return tile, nil
	}

	remaining := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	if remaining <= params.WarnDays {
		tile.Status = coreModels.WarningStatus
	} else {
		tile.Status = coreModels.SuccessStatus
	}

	tile.Message = buildMessage(params.Display, cert, remaining)

	tile.WithMetrics(coreModels.RawUnit)

	tile.Metrics.Values = []string{
		fmt.Sprintf("Expires in %d days", remaining),
	}

	return tile, nil
}

func buildMessage(display string, cert *api.Certificate, remaining int) string {
	if display == "" || display == "full" {
		return fmt.Sprintf(
			"NotBefore: %s / NotAfter: %s / Issuer: %s / CN: %s",
			cert.NotBefore.Format(time.RFC3339),
			cert.NotAfter.Format(time.RFC3339),
			cert.Issuer,
			cert.Subject,
		)
	}

	var parts []string
	for _, field := range strings.Split(display, ",") {
		switch strings.ToLower(strings.TrimSpace(field)) {
		case "remaining":
			parts = append(parts, fmt.Sprintf("%d", remaining))
		case "notbefore":
			parts = append(parts, cert.NotBefore.Format(time.RFC3339))
		case "notafter":
			parts = append(parts, cert.NotAfter.Format(time.RFC3339))
		case "issuer":
			parts = append(parts, cert.Issuer)
		case "subject":
			parts = append(parts, cert.Subject)
		}
	}

	return strings.Join(parts, " / ")
}
