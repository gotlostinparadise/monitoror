//go:build !faker

package dns

import (
    "github.com/monitoror/monitoror/api/config/versions"
    coreModels "github.com/monitoror/monitoror/models"
    "github.com/monitoror/monitoror/monitorables/dns/api"
    dnsDelivery "github.com/monitoror/monitoror/monitorables/dns/api/delivery/http"
    dnsModels "github.com/monitoror/monitoror/monitorables/dns/api/models"
    dnsRepository "github.com/monitoror/monitoror/monitorables/dns/api/repository"
    dnsUsecase "github.com/monitoror/monitoror/monitorables/dns/api/usecase"
    "github.com/monitoror/monitoror/registry"
    "github.com/monitoror/monitoror/store"
)

type Monitorable struct {
    store *store.Store

    dnsTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
    m := &Monitorable{store: store}
    m.dnsTileEnabler = store.Registry.RegisterTile(api.DNSTileType, versions.MinimalVersion, m.GetVariantsNames())
    return m
}

func (m *Monitorable) GetDisplayName() string { return "DNS" }

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName { return []coreModels.VariantName{coreModels.DefaultVariantName} }

func (m *Monitorable) Validate(_ coreModels.VariantName) (bool, []error) { return true, nil }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
    repository := dnsRepository.NewDNSRepository()
    usecase := dnsUsecase.NewDNSUsecase(repository)
    delivery := dnsDelivery.NewDNSDelivery(usecase)

    routeGroup := m.store.MonitorableRouter.Group("/dns", variantName)
    route := routeGroup.GET("/dns", delivery.GetDNS)

    m.dnsTileEnabler.Enable(variantName, &dnsModels.DNSParams{}, route.Path)
}
