package config

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfiguration
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port string
}

type DatabaseConfiguration struct {
	PgName string
}
