package repository

import (
    "net"

    "github.com/monitoror/monitoror/monitorables/dns/api"
)

type dnsRepository struct{}

func NewDNSRepository() api.Repository { return &dnsRepository{} }

func (r *dnsRepository) Lookup(recordType, name string) ([]string, error) {
    switch recordType {
    case "A":
        ips, err := net.LookupIP(name)
        if err != nil {
            return nil, err
        }
        var results []string
        for _, ip := range ips {
            if ip.To4() != nil {
                results = append(results, ip.String())
            }
        }
        return results, nil
    case "AAAA":
        ips, err := net.LookupIP(name)
        if err != nil {
            return nil, err
        }
        var results []string
        for _, ip := range ips {
            if ip.To16() != nil && ip.To4() == nil {
                results = append(results, ip.String())
            }
        }
        return results, nil
    case "CNAME":
        cname, err := net.LookupCNAME(name)
        if err != nil {
            return nil, err
        }
        return []string{cname}, nil
    case "TXT":
        txt, err := net.LookupTXT(name)
        if err != nil {
            return nil, err
        }
        return txt, nil
    default:
        return nil, net.UnknownNetworkError(recordType)
    }
}
