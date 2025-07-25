//go:build !faker

package rtsp

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	rtspDelivery "github.com/monitoror/monitoror/monitorables/rtsp/api/delivery/http"
	rtspModels "github.com/monitoror/monitoror/monitorables/rtsp/api/models"
	rtspRepository "github.com/monitoror/monitoror/monitorables/rtsp/api/repository"
	rtspUsecase "github.com/monitoror/monitoror/monitorables/rtsp/api/usecase"
	rtspConfig "github.com/monitoror/monitoror/monitorables/rtsp/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*rtspConfig.RTSP

	rtspTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*rtspConfig.RTSP)

	pkgMonitorable.LoadConfig(&m.config, rtspConfig.Default)

	m.rtspTileEnabler = store.Registry.RegisterTile(api.RTSPTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "RTSP" }

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

	repository := rtspRepository.NewRTSPRepository(conf)
	usecase := rtspUsecase.NewRTSPUsecase(repository)
	delivery := rtspDelivery.NewRTSPDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/rtsp", variantName)
	route := routeGroup.GET("/rtsp", delivery.GetRTSP)

	m.rtspTileEnabler.Enable(variantName, &rtspModels.RTSPParams{}, route.Path)
}
