//go:build faker

package command

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/command/api"
	commandDelivery "github.com/monitoror/monitoror/monitorables/command/api/delivery/http"
	commandModels "github.com/monitoror/monitoror/monitorables/command/api/models"
	commandUsecase "github.com/monitoror/monitoror/monitorables/command/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	commandTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{store: store}

	m.commandTileEnabler = store.Registry.RegisterTile(api.CommandTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Command" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := commandUsecase.NewCommandUsecase(nil)
	delivery := commandDelivery.NewCommandDelivery(usecase)

	routeGroup := m.store.MonitorableRouter.Group("/command", variantName)
	route := routeGroup.GET("/status", delivery.GetCommandStatus)

	m.commandTileEnabler.Enable(variantName, &commandModels.CommandParams{}, route.Path)
}
