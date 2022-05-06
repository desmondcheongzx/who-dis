package pkg

import (
	"net"
)

type DNSClient struct {
	cache map[string]net.IP
}

func NewDNSClient() *DNSClient {
	return &DNSClient{
		cache: make(map[string]net.IP),
	}
}

func (client *DNSClient) Query(dn string) net.IP {
	return nil
}
