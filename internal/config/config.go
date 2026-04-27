package config

import (
	"flag"
	"os"

	"net"
	"net/url"
	"strings"
)

// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
// Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).

const (
	defaultServerAddress = ":8080"
	defaultBaseAddress   = "http://localhost:8080"

	flagServerAddress = "a"
	flagBaseAddress   = "b"
	envBaseAddress    = "BASE_URL"
	envServerAddress  = "SERVER_ADDRESS"
	descServerAddress = "адрес запуска HTTP-сервера"
	descBaseAddress   = "базовый адрес результирующего сокращённого URL"

	defaultLogLevel = "Debug"
	flagLogLevel    = "l"
	envLogLevel     = "LOG_LEVEL"
	descLogLevel    = "уровень логирования"
)

var DefaultConfig = &Config{
	serverAddress: defaultServerAddress,
	baseAddress:   defaultBaseAddress,
	logLevel:      defaultLogLevel,
}

type Config struct {
	serverAddress string
	baseAddress   string
	logLevel      string
}

func New() *Config {

	serverAddress := setAddress(envServerAddress, flagServerAddress, defaultServerAddress, descServerAddress)
	baseAddress := setAddress(envBaseAddress, flagBaseAddress, defaultBaseAddress, descBaseAddress)
	logLevel := setAddress(envLogLevel, flagLogLevel, defaultLogLevel, descLogLevel)

	flag.Parse()

	return &Config{
		serverAddress: validateServerAddress(*serverAddress, defaultServerAddress),
		baseAddress:   validateBaseAddress(*baseAddress, defaultBaseAddress),
		logLevel:      *logLevel,
	}
}

func setAddress(envAddress, flagName, defaultAddress, description string) *string {
	flagaddress := flag.String(flagName, defaultAddress, description)
	if address, ok := os.LookupEnv(envAddress); ok && address != "" {
		return &address
	}
	return flagaddress
}

func validateServerAddress(address, defaultAddress string) string {
	addrList := strings.Split(address, ":")
	if len(addrList) < 1 || len(addrList) > 2 || len(addrList) == 1 && addrList[0] == "" {
		return defaultAddress
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

func (c *Config) GetLogLevel() string {
	return c.logLevel
}
