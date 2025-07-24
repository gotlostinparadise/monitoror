//go:build faker

package usecase

import (
	"fmt"
	"strings"
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
	tile.Label = fmt.Sprintf("%s:%d", params.Domain, params.GetPort())

	status := su.computeStatus(params)
	tile.Status = nonempty.Struct(params.Status, status).(coreModels.TileStatus)
	remaining := su.computeRemainingDays(params.Domain)
	cert := &api.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour * 24 * 90),
		Issuer:    "issuer",
		Subject:   "subject",
	}
	tile.Message = buildMessage(params.Display, cert, remaining)
	tile.Metrics.Values = []string{"", "", "issuer", "subject"}
	return tile, nil
}

func (su *sslUsecase) computeStatus(params *models.SSLParams) coreModels.TileStatus {
	value, ok := su.timeRefByHost[params.Domain]
	if !ok {
		su.timeRefByHost[params.Domain] = faker.GetRefTime()
	}
	return faker.ComputeStatus(value, availableStatuses)
}

func (su *sslUsecase) computeRemainingDays(hostname string) int {
	value, ok := su.timeRefByHost[hostname]
	if !ok {
		su.timeRefByHost[hostname] = faker.GetRefTime()
		value = su.timeRefByHost[hostname]
	}
	duration := faker.ComputeDuration(value, time.Hour*24*90)
	remaining := 90 - int(duration.Hours()/24)
	if remaining == 0 {
		remaining = 90
	}
	return remaining
}

func buildMessage(display string, cert *api.Certificate, remaining int) string {
	if display == "" || display == "full" {
		return fmt.Sprintf(
			"Expires in %d days / NotBefore: %s / NotAfter: %s / Issuer: %s",
			remaining,
			cert.NotBefore.Format(time.RFC3339),
			cert.NotAfter.Format(time.RFC3339),
			cert.Issuer,
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
