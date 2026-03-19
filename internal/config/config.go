package config

import (
	"flag"
	"net"
	"net/url"
	"strings"
)

// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
// Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).

const (
	defaultServerAddress = ""
	defaultServerPort    = "8080"
	defaultBaseAddress   = "http://localhost:8080"
)

type Config struct {
	serverAddress string
	baseAddress   string
}

func New() *Config {
	res := &Config{}
	flag.StringVar(&res.serverAddress, "a", ":8080", "адрес HTTP-сервера")
	flag.StringVar(&res.baseAddress, "b", defaultBaseAddress, "базовый адрес сокращённого URL")

	flag.Parse()

	res.serverAddress = validateServerAddress(res.serverAddress, defaultServerAddress, defaultServerPort)
	res.baseAddress = validateBaseAddress(res.baseAddress, defaultBaseAddress)

	return res
}

func validateServerAddress(address, defaultAddress, defaultPort string) string {
	addrList := strings.Split(address, ":")
	if len(addrList) < 1 || len(addrList) > 2 {
		return net.JoinHostPort(defaultAddress, defaultPort)
	}
	if addrList[0] == "" {
		addrList[0] = defaultAddress
	}
	if len(addrList) < 2 || addrList[1] == "" {
		return addrList[0]
	}
	return net.JoinHostPort(addrList[0], addrList[1])
}

func validateBaseAddress(address, defaultAddress string) string {
	u, err := url.Parse(address)
	if err != nil {
		return defaultAddress
	}
	if u.Scheme == "" || u.Host == "" {
		return defaultAddress
	}

	return u.String()
}

func (c *Config) GetBaseAddress() string {
	return c.baseAddress
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
}
