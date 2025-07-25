//go:build faker

package whois

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api"
	whoisDelivery "github.com/monitoror/monitoror/monitorables/whois/api/delivery/http"
	whoisModels "github.com/monitoror/monitoror/monitorables/whois/api/models"
	whoisUsecase "github.com/monitoror/monitoror/monitorables/whois/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	whoisTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.whoisTileEnabler = store.Registry.RegisterTile(api.WHOISTileType, versions.MinimalVersion, m.GetVariantsNames())
	return m
}

func (m *Monitorable) GetDisplayName() string { return "WHOIS" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := whoisUsecase.NewWHOISUsecase()
	delivery := whoisDelivery.NewWHOISDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/whois", variantName)
	route := routeGroup.GET("/domain", delivery.GetWHOIS)

	m.whoisTileEnabler.Enable(variantName, &whoisModels.WHOISParams{}, route.Path)
}
