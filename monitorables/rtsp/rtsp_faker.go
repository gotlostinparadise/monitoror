//go:build faker

package rtsp

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	rtspDelivery "github.com/monitoror/monitoror/monitorables/rtsp/api/delivery/http"
	rtspModels "github.com/monitoror/monitoror/monitorables/rtsp/api/models"
	rtspUsecase "github.com/monitoror/monitoror/monitorables/rtsp/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	rtspTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	m.rtspTileEnabler = store.Registry.RegisterTile(api.RTSPTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "RTSP" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := rtspUsecase.NewRTSPUsecase(nil)
	delivery := rtspDelivery.NewRTSPDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/rtsp", variantName)
	route := routeGroup.GET("/rtsp", delivery.GetRTSP)

	m.rtspTileEnabler.Enable(variantName, &rtspModels.RTSPParams{}, route.Path)
}
