package config

// Config представляет конфигурацию приложения.
type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
	Clients ClientsConfig `toml:"clients"`
}

// GlobalConfig представляет глобальные настройки.
type GlobalConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"dev", "stage", "prod"}.
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

// LogConfig представляет настройки логирования.
type LogConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"debug", "info", "warn", "error"}.
	Level          string `toml:"level" validate:"required,oneof=debug info warn error"`
	ProductionMode bool   `toml:"production_mode"`
}

// ServersConfig представляет настройки серверов.
type ServersConfig struct {
	Debug  DebugServerConfig  `toml:"debug"`
	Client ClientServerConfig `toml:"client"`
}

// DebugServerConfig представляет настройки отладочного сервера.
type DebugServerConfig struct {
	// добавляем валидацию: обязательное поле, значение должно быть в формате "host:port".
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

// SentryConfig представляет настройки Sentry.
type SentryConfig struct {
	// DSN - URL для отправки отчетов в Sentry.
	// добавляем валидацию: значение должно быть в формате URL, не работает если поле пустое
	DSN string `toml:"dsn" validate:"omitempty,url"`
}

// ClientServerConfig представляет настройки клиентского сервера.
type ClientServerConfig struct {
	Addr         string   `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins []string `toml:"allow_origins" validate:"required,dive,uri"`
}

// ClientsConfig представляет настройки для Keycloak.
type ClientsConfig struct {
	Keycloak KeycloakConfig `toml:"keycloak"`
}

// KeycloakConfig представляет настройки для Keycloak.
type KeycloakConfig struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}
