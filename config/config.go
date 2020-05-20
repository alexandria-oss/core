package config

import (
	"context"

	"github.com/spf13/viper"
)

// Kernel Alexandria kernel configuration struct
// Generates required OS env variables
type Kernel struct {
	Transport transport
	Tracing   tracing

	EventBus eventBus

	Docstore docstore
	DBMS     dbms
	InMemory inMemory

	AWS aws

	Auth auth

	Version string
	Service string
}

func init() {
	// Service info
	viper.SetDefault("alexandria.info.version", "0.1.0")
	viper.SetDefault("alexandria.info.service", "example-service")
}

// NewKernel Generate a global configuration from alexandria-config.yml file
func NewKernel(ctx context.Context) (*Kernel, error) {
	// Context is required to use gocloud.dev functions

	kernelConfig := new(Kernel)

	// Init config
	viper.SetConfigName("alexandria-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("/etc/alexandria/")
	viper.AddConfigPath("$HOME/.alexandria")
	viper.AddConfigPath(".")

	// Open config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			_ = viper.SafeWriteConfig()
		}

		// Config file was found but another error was produced, use default values
	}

	// Map kernel configuration
	kernelConfig.Transport = newTransportConfig()
	kernelConfig.Tracing = newTracingConfig()
	kernelConfig.EventBus = newEventBusConfig()
	kernelConfig.Docstore = newDocstoreConfig()
	kernelConfig.DBMS = newDBMSConfig()
	kernelConfig.InMemory = newInMemoryConfig()
	kernelConfig.AWS = newAWSConfig()
	kernelConfig.Auth = newAuthConfig()

	kernelConfig.Version = viper.GetString("alexandria.info.version")
	kernelConfig.Service = viper.GetString("alexandria.info.service")

	// Prefer AWS KMS/Hashicorp Vault/Key Parameter Store over local, replace default or local config
	// TODO: Implement Hashicorp Vault or AWS KMS key/secret fetching

	return kernelConfig, nil
}
