package config

// Config is an application config
// Should be used only in main packages for config parsing and dependency initialization.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig

	Log LogConfig
}

type ServerConfig struct {
	Port    int
	TLSKey  string
	TLSCert string
}

type DatabaseConfig struct {
	URI       string
	MinConns  int32
	MaxConns  int32
	Migration MigrationConfig
}

type MigrationConfig struct {
	Enabled      bool
	Path         string
	VersionTable string
}

type LogConfig struct {
	Development bool
	Verbosity   int8
}
