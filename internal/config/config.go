/*
Package config конфигурация

	Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
	Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).
	Флаг -l отвечает за уровень логирования (значение по умолчанию: "Info")
	Флаг -f путь до файла, куда сохраняются данные в формате JSON (значение по умолчанию "./storage.json")
*/
package config

import (
	"flag"
	"net"
	"net/url"
	"os"
	"strings"
)

// ShortLen — длина генерируемого короткого URL в символах.
const ShortLen = 6

const (
	defaultServerAddress = ":8080"
	flagServerAddress    = "a"
	envServerAddress     = "SERVER_ADDRESS"
	descServerAddress    = "адрес запуска HTTP-сервера"

	defaultBaseAddress = "http://localhost:8080"
	flagBaseAddress    = "b"
	envBaseAddress     = "BASE_URL"
	descBaseAddress    = "базовый адрес результирующего сокращённого URL"

	defaultLogLevel = "Info"
	flagLogLevel    = "l"
	envLogLevel     = "LOG_LEVEL"
	descLogLevel    = "уровень логирования"

	defaultFileName = ""
	flagFileName    = "f"
	envFileName     = "FILE_STORAGE_PATH"
	descFileName    = "файл для хранения сокращенных адресов"

	defaultConnAddress = ""
	flagConnAddress    = "d"
	envConnAddress     = "DATABASE_DSN"
	descConnAddress    = "строка с адресом подключения к БД"

	defaultAuditFile = ""
	flagAuditFile    = "audit-file"
	envAuditFile     = "AUDIT_FILE"
	descAuditFile    = "путь к файлу-приёмнику, в который сохраняются логи аудита"

	defaultAuditAddress = ""
	flagAuditAddress    = "audit-url"
	envAuditAddress     = "AUDIT_URL"
	descAuditAddress    = "полный URL удаленного сервера-приёмника, куда отправляются логи аудита"
)

// DefaultConfig — конфигурация по умолчанию.
var DefaultConfig = &Config{
	serverAddress: defaultServerAddress,
	baseAddress:   defaultBaseAddress,
	logLevel:      defaultLogLevel,
	fileName:      defaultFileName,
	connAddress:   defaultConnAddress,
}

// Config содержит параметры конфигурации сервера.
type Config struct {
	serverAddress string
	baseAddress   string
	logLevel      string
	fileName      string
	connAddress   string
	auditFile     string
	auditAddress  string
}

// New создаёт и возвращает конфигурацию, читая значения из флагов и переменных окружения.
func New() *Config {

	serverAddress := setAddress(envServerAddress, flagServerAddress, defaultServerAddress, descServerAddress)
	baseAddress := setAddress(envBaseAddress, flagBaseAddress, defaultBaseAddress, descBaseAddress)
	logLevel := setAddress(envLogLevel, flagLogLevel, defaultLogLevel, descLogLevel)
	fileName := setAddress(envFileName, flagFileName, defaultFileName, descFileName)
	connAddress := setAddress(envConnAddress, flagConnAddress, defaultConnAddress, descConnAddress)
	auditFile := setAddress(envAuditFile, flagAuditFile, defaultAuditFile, descAuditFile)
	auditAddress := setAddress(envAuditAddress, flagAuditAddress, defaultAuditAddress, descAuditAddress)

	flag.Parse()

	return &Config{
		serverAddress: validateServerAddress(*serverAddress, defaultServerAddress),
		baseAddress:   validateBaseAddress(*baseAddress, defaultBaseAddress),
		logLevel:      *logLevel,
		fileName:      *fileName,
		connAddress:   *connAddress,
		auditFile:     *auditFile,
		auditAddress:  validateBaseAddress(*auditAddress, defaultAuditAddress),
	}
}

// GetBaseAddress возвращает базовый адрес результирующего сокращённого URL.
func (c *Config) GetBaseAddress() string {
	return c.baseAddress
}

// GetServerAddress возвращает адрес запуска HTTP-сервера.
func (c *Config) GetServerAddress() string {
	return c.serverAddress
}

// GetLogLevel возвращает уровень логирования.
func (c *Config) GetLogLevel() string {
	return c.logLevel
}

// GetFileName возвращает путь к файлу для хранения сокращённых адресов.
func (c *Config) GetFileName() string {
	return c.fileName
}

// GetConnAddress возвращает строку подключения к базе данных.
func (c *Config) GetConnAddress() string {
	return c.connAddress
}

// GetAuditFile возвращает путь к файлу-приёмнику логов аудита.
func (c *Config) GetAuditFile() string {
	return c.auditFile
}

// GetAuditAddress возвращает URL удалённого сервера-приёмника логов аудита.
func (c *Config) GetAuditAddress() string {
	return c.auditAddress
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
