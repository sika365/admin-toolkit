package config

import "github.com/alifakhimi/simple-service-go"

type ServiceConfig struct {
	*simple.Config `mapstructure:",squash"`
	// custom meta
	// Meta *Meta `json:"meta,omitempty" mapstructure:"meta"`
	Meta map[string]any `json:"meta,omitempty" mapstructure:"meta"`
}

type Meta struct {
	MetaValue `mapstructure:",squash"`
	Mock      MetaValue `json:"mock,omitempty" mapstructure:"mock"`
}

type MetaValue struct {
}

var (
	conf = ServiceConfig{}
)

func Config() *ServiceConfig {
	return &conf
}
