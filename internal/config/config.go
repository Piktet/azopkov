/*
Package config конфигурация

	Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
	Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).
	Флаг -l отвечает за уровень логирования (значение по умолчанию: "Info")
	Флаг -f путь до файла, куда сохраняются данные в формате JSON (значение по умолчанию "./storage.json")
*/
package config

import (
	"encoding/json"
	"flag"
	"os"
	"sync"
)

// ShortLen — длина генерируемого короткого URL в символах.
const ShortLen = 6

type sourceType int

const (
	_ sourceType = iota
	sourceEnv
	sourceFlag
	sourceConfig
	sourceDefault
)

type valueTypeType int

const (
	_ valueTypeType = iota
	valueString
	valueBool
)

type ConfigValue struct {
	defaultValue any
	flagName     []string
	envName      string
	description  string

	value     any
	valueType valueTypeType

	source sourceType
}

type configType string

const (
	configServerAddress = "server_address"
	configBaseAddress   = "base_url"
	configLogLevel      = "log_level"
	configFileName      = "file_storage_path"
	configConnAddress   = "database_dsn"
	configAuditFile     = "audit_file"
	configAuditAddress  = "audit_url"
	configEnableHTTPS   = "enable_https"
	configConfigFile    = "config_file"
)

const (
	defaultServerAddress = ":8080"
	defaultBaseAddress   = "http://localhost:8080"
	defaultLogLevel      = "Info"
)

var (
	configData = map[configType]*ConfigValue{
		configServerAddress: {defaultValue: defaultServerAddress, flagName: []string{"a"}, envName: "SERVER_ADDRESS", description: "", valueType: valueString},
		configBaseAddress:   {defaultValue: defaultBaseAddress, flagName: []string{"b"}, envName: "BASE_URL", description: "", source: sourceDefault, valueType: valueString},
		configLogLevel:      {defaultValue: defaultLogLevel, flagName: []string{"l"}, envName: "LOG_LEVEL", description: "", valueType: valueString},
		configFileName:      {defaultValue: "", flagName: []string{"f"}, envName: "FILE_STORAGE_PATH", description: "", valueType: valueString},
		configConnAddress:   {defaultValue: "", flagName: []string{"d"}, envName: "DATABASE_DSN", description: "", valueType: valueString},
		configAuditFile:     {defaultValue: "", flagName: []string{"audit-file"}, envName: "AUDIT_FILE", description: "", valueType: valueString},
		configAuditAddress:  {defaultValue: "", flagName: []string{"audit-url"}, envName: "AUDIT_URL", description: "", valueType: valueString},
		configEnableHTTPS:   {defaultValue: false, flagName: []string{"s"}, envName: "ENABLE_HTTPS", description: "", valueType: valueBool},
		configConfigFile:    {defaultValue: "", flagName: []string{"c", "config"}, envName: "CONFIG", description: "файл конфигурации", valueType: valueString},
	}
)

// Config содержит параметры конфигурации сервера.
type Config struct {
	list map[configType]*ConfigValue
}

var fn sync.Once

// New создаёт и возвращает конфигурацию, читая значения из флагов и переменных окружения.
func New() *Config {
	fn.Do(func() {

		for _, v := range configData {
			v.value = v.defaultValue
			v.source = sourceDefault
			switch v.valueType {

			case valueString:
				for _, f := range v.flagName {
					_ = flag.String(f, v.defaultValue.(string), v.description)
				}
			case valueBool:
				for _, f := range v.flagName {
					_ = flag.Bool(f, v.defaultValue.(bool), v.description)
				}
			}
			if respEnv, ok := os.LookupEnv(v.envName); ok && respEnv != "" {
				v.value = respEnv
				v.source = sourceEnv
			}
		}

		flag.Parse()

		flag.Visit(func(flagValue *flag.Flag) {
			if v, ok := configData[configType(flagValue.Name)]; ok {
				v.source = sourceFlag
				switch v.valueType {
				case valueString:
					v.value = flagValue.Value.String()
				case valueBool:
					v.value = flagValue.Value.String() == "true"
				}
			}
		})

		if confFile, ok := configData[configConfigFile]; ok {
			if data, err := os.ReadFile(confFile.value.(string)); err == nil {
				val := make(map[string]any)
				if err = json.Unmarshal(data, &val); err == nil {
					for k, v := range val {
						if item, ok := configData[configType(k)]; ok {
							if item.source == sourceDefault {
								item.source = sourceConfig
								switch item.valueType {
								case valueString:
									item.value = v.(string)
								case valueBool:
									item.value = v.(bool)
								}
							}
						}
					}
				}
			}
		}

	})
	return &Config{list: configData}
}

func (c *Config) getString(name configType) string {
	if v, ok := c.list[name]; ok && v != nil && v.valueType == valueString && v.value != nil {
		return v.value.(string)
	}
	return ""
}

func (c *Config) getBool(name configType) bool {
	if v, ok := c.list[name]; ok && v != nil && v.valueType == valueBool && v.value != nil {
		return v.value.(bool)
	}
	return false

}

// GetBaseAddress возвращает базовый адрес результирующего сокращённого URL.
func (c *Config) GetBaseAddress() string {
	return c.getString(configBaseAddress)
}

// GetServerAddress возвращает адрес запуска HTTP-сервера.
func (c *Config) GetServerAddress() string {
	return c.getString(configServerAddress)
}

// GetLogLevel возвращает уровень логирования.
func (c *Config) GetLogLevel() string {
	return c.getString(configLogLevel)
}

// GetFileName возвращает путь к файлу для хранения сокращённых адресов.
func (c *Config) GetFileName() string {
	return c.getString(configFileName)
}

// GetConnAddress возвращает строку подключения к базе данных.
func (c *Config) GetConnAddress() string {
	return c.getString(configConnAddress)
}

// GetAuditFile возвращает путь к файлу-приёмнику логов аудита.
func (c *Config) GetAuditFile() string {
	return c.getString(configAuditFile)
}

// GetAuditAddress возвращает URL удалённого сервера-приёмника логов аудита.
func (c *Config) GetAuditAddress() string {
	return c.getString(configAuditAddress)
}

// IsEnableHTTPS использовать HTTPS
func (c *Config) IsEnableHTTPS() bool {
	return c.getBool(configEnableHTTPS)
}
