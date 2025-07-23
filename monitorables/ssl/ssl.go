//go:build !faker

package ssl

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	sslDelivery "github.com/monitoror/monitoror/monitorables/ssl/api/delivery/http"
	sslModels "github.com/monitoror/monitoror/monitorables/ssl/api/models"
	sslRepository "github.com/monitoror/monitoror/monitorables/ssl/api/repository"
	sslUsecase "github.com/monitoror/monitoror/monitorables/ssl/api/usecase"
	sslConfig "github.com/monitoror/monitoror/monitorables/ssl/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*sslConfig.SSL

	sslTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*sslConfig.SSL)

	pkgMonitorable.LoadConfig(&m.config, sslConfig.Default)

	m.sslTileEnabler = store.Registry.RegisterTile(api.SSLTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "SSL" }

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	if errors := pkgMonitorable.ValidateConfig(conf, variantName); errors != nil {
		return false, errors
	}

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := sslRepository.NewSSLRepository(conf)
	usecase := sslUsecase.NewSSLUsecase(repository)
	delivery := sslDelivery.NewSSLDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/ssl", variantName)
	route := routeGroup.GET("/cert", delivery.GetSSL)

	m.sslTileEnabler.Enable(variantName, &sslModels.SSLParams{}, route.Path)
}
