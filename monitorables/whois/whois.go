//go:build !faker

package whois

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api"
	whoisDelivery "github.com/monitoror/monitoror/monitorables/whois/api/delivery/http"
	whoisModels "github.com/monitoror/monitoror/monitorables/whois/api/models"
	whoisRepository "github.com/monitoror/monitoror/monitorables/whois/api/repository"
	whoisUsecase "github.com/monitoror/monitoror/monitorables/whois/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*struct{}

	whoisTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*struct{})

	pkgMonitorable.LoadConfig(&m.config, &struct{}{})
	if len(m.config) == 0 {
		m.config[coreModels.DefaultVariantName] = &struct{}{}
	}

	m.whoisTileEnabler = store.Registry.RegisterTile(api.WHOISTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "WHOIS" }

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	if errors := pkgMonitorable.ValidateConfig(m.config[variantName], variantName); errors != nil {
		return false, errors
	}
	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	repository := whoisRepository.NewWHOISRepository()
	usecase := whoisUsecase.NewWHOISUsecase(repository)
	delivery := whoisDelivery.NewWHOISDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/whois", variantName)
	route := routeGroup.GET("/domain", delivery.GetWHOIS)

	m.whoisTileEnabler.Enable(variantName, &whoisModels.WHOISParams{}, route.Path)
}
