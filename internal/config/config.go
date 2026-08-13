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

	defaultEnableHTTPS = false
	flagEnableHTTPS    = "s"
	envEnableHTTPS     = "ENABLE_HTTPS"
	descEnableHTTPS    = "использовать HTTPS"
)

// DefaultConfig — конфигурация по умолчанию.
var DefaultConfig = &Config{
	serverAddress: defaultServerAddress,
	baseAddress:   defaultBaseAddress,
	logLevel:      defaultLogLevel,
	fileName:      defaultFileName,
	connAddress:   defaultConnAddress,
	auditFile:     defaultAuditFile,
	auditAddress:  defaultAuditAddress,
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
	enableHTTPS   bool
}

// New создаёт и возвращает конфигурацию, читая значения из флагов и переменных окружения.
func New() *Config {

	serverAddress := setParamString(envServerAddress, flagServerAddress, defaultServerAddress, descServerAddress)
	baseAddress := setParamString(envBaseAddress, flagBaseAddress, defaultBaseAddress, descBaseAddress)
	logLevel := setParamString(envLogLevel, flagLogLevel, defaultLogLevel, descLogLevel)
	fileName := setParamString(envFileName, flagFileName, defaultFileName, descFileName)
	connAddress := setParamString(envConnAddress, flagConnAddress, defaultConnAddress, descConnAddress)
	auditFile := setParamString(envAuditFile, flagAuditFile, defaultAuditFile, descAuditFile)
	auditAddress := setParamString(envAuditAddress, flagAuditAddress, defaultAuditAddress, descAuditAddress)
	enableHTTPS := setParamBool(envEnableHTTPS, flagEnableHTTPS, defaultEnableHTTPS, descEnableHTTPS)
	flag.Parse()

	return &Config{
		serverAddress: validateServerAddress(*serverAddress, defaultServerAddress),
		baseAddress:   validateBaseAddress(*baseAddress, defaultBaseAddress),
		logLevel:      *logLevel,
		fileName:      *fileName,
		connAddress:   *connAddress,
		auditFile:     *auditFile,
		auditAddress:  validateBaseAddress(*auditAddress, defaultAuditAddress),
		enableHTTPS:   *enableHTTPS,
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

func setParamString(env, flagName, defaultAddress, description string) *string {
	respFlag := flag.String(flagName, defaultAddress, description)
	if respEnv, ok := os.LookupEnv(env); ok && respEnv != "" {
		return &respEnv
	}
	return respFlag
}

func setParamBool(env, flagName string, defaultParam bool, description string) *bool {
	respFflag := flag.Bool(flagName, defaultParam, description)
	if x, ok := os.LookupEnv(env); ok && x != "" {
		respEnv := true
		return &respEnv
	}
	return respFflag
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

// IsEnableHTTPS использовать HTTPS
func (c *Config) IsEnableHTTPS() bool {
	return c.enableHTTPS
}
