package config

// Config представляет конфигурацию приложения.
type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
}

// GlobalConfig представляет глобальные настройки.
type GlobalConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"dev", "stage", "prod"}.
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

// LogConfig представляет настройки логирования.
type LogConfig struct {
	// добавляем валидацию: обязательное поле, значения из {"debug", "info", "warn", "error"}.
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

// ServersConfig представляет настройки серверов.
type ServersConfig struct {
	Debug DebugServerConfig `toml:"debug"`
}

// DebugServerConfig представляет настройки отладочного сервера.
type DebugServerConfig struct {
	// добавляем валидацию: обязательное поле, значение должно быть в формате "host:port".
	Addr string `toml:"addr" validate:"required,hostname_port"`
}
