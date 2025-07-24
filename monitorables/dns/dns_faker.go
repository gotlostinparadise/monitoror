//go:build faker

package dns

import (
    "github.com/monitoror/monitoror/api/config/versions"
    "github.com/monitoror/monitoror/internal/pkg/monitorable"
    coreModels "github.com/monitoror/monitoror/models"
    "github.com/monitoror/monitoror/monitorables/dns/api"
    dnsDelivery "github.com/monitoror/monitoror/monitorables/dns/api/delivery/http"
    dnsModels "github.com/monitoror/monitoror/monitorables/dns/api/models"
    dnsUsecase "github.com/monitoror/monitoror/monitorables/dns/api/usecase"
    "github.com/monitoror/monitoror/registry"
    "github.com/monitoror/monitoror/store"
)

type Monitorable struct {
    monitorable.DefaultMonitorableFaker

    store *store.Store

    dnsTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
    m := &Monitorable{store: store}
    m.dnsTileEnabler = store.Registry.RegisterTile(api.DNSTileType, versions.MinimalVersion, m.GetVariantsNames())
    return m
}

func (m *Monitorable) GetDisplayName() string { return "DNS" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
    usecase := dnsUsecase.NewDNSUsecase()
    delivery := dnsDelivery.NewDNSDelivery(usecase)

    routeGroup := m.store.MonitorableRouter.Group("/dns", variantName)
    route := routeGroup.GET("/dns", delivery.GetDNS)

    m.dnsTileEnabler.Enable(variantName, &dnsModels.DNSParams{}, route.Path)
}
