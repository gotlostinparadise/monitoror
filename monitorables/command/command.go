//go:build !faker

package command

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	commandDelivery "github.com/monitoror/monitoror/monitorables/command/api/delivery/http"
	commandModels "github.com/monitoror/monitoror/monitorables/command/api/models"
	commandRepository "github.com/monitoror/monitoror/monitorables/command/api/repository"
	commandUsecase "github.com/monitoror/monitoror/monitorables/command/api/usecase"
	commandConfig "github.com/monitoror/monitoror/monitorables/command/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*commandConfig.Command

	commandTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*commandConfig.Command)

	pkgMonitorable.LoadConfig(&m.config, commandConfig.Default)

	m.commandTileEnabler = store.Registry.RegisterTile(api.CommandTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Command" }

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

	repository := commandRepository.NewCommandRepository(conf)
	usecase := commandUsecase.NewCommandUsecase(repository)
	delivery := commandDelivery.NewCommandDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/command", variantName)
	route := routeGroup.GET("/status", delivery.GetCommandStatus)

	m.commandTileEnabler.Enable(variantName, &commandModels.CommandParams{}, route.Path)
}
