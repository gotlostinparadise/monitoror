//go:build faker

package ssl

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ssl/api"
	sslDelivery "github.com/monitoror/monitoror/monitorables/ssl/api/delivery/http"
	sslModels "github.com/monitoror/monitoror/monitorables/ssl/api/models"
	sslUsecase "github.com/monitoror/monitoror/monitorables/ssl/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	sslTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	m.sslTileEnabler = store.Registry.RegisterTile(api.SSLTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "SSL" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := sslUsecase.NewSSLUsecase(nil)
	delivery := sslDelivery.NewSSLDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/ssl", variantName)
	route := routeGroup.GET("/cert", delivery.GetSSL)

	m.sslTileEnabler.Enable(variantName, &sslModels.SSLParams{}, route.Path)
}
