package config

import (
	"fmt"
	"strings"
)

// AddressMap is a mapping of service -> service address
type AddressMap struct {
	Hook       ServiceAddress
	Data       ServiceAddress
	Queue      ServiceAddress
	UI         ServiceAddress
	Asset      ServiceAddress
	Repository ServiceAddress
	Auth       ServiceAddress
	Log        ServiceAddress
}

// DefaultServices is the standard array of service mappings when unconfigured.
var DefaultServices = AddressMap{
	Hook:       ServiceAddress{Port: 2020},
	Data:       ServiceAddress{Port: 6000},
	Queue:      ServiceAddress{Port: 6001},
	Asset:      ServiceAddress{Port: 6002},
	Repository: ServiceAddress{Port: 6003},
	Auth:       ServiceAddress{Port: 6004},
	Log:        ServiceAddress{Port: 6005},
	UI:         ServiceAddress{Port: 6010, HTTP: true},
}

// ServiceAddress is a well-formed address for service connections
type ServiceAddress struct {
	Hostname string
	Port     uint
	HTTP     bool
	TLS      bool
}

func (sa ServiceAddress) String() string {
	hostname := sa.Hostname
	if hostname == "" {
		hostname = "localhost"
	}

	str := strings.Join([]string{hostname, fmt.Sprintf("%v", sa.Port)}, ":")
	if sa.HTTP {
		if sa.TLS {
			str = "https://" + str
		} else {
			str = "http://" + str
		}
	}

	return str
}
