package config

import (
	"github.com/spf13/viper"
	"os"
)

type eventBus struct {
	KafkaBrokers []string
}

func init() {
	viper.SetDefault("alexandria.eventbus.kafka.brokers", []string{"0.0.0.0:9092"})
}

func newEventBusConfig() eventBus {
	cfg := eventBus{
		KafkaBrokers: viper.GetStringSlice("alexandria.eventbus.kafka.brokers"),
	}

	// Start up required kafka env
	_ = os.Setenv("KAFKA_BROKERS", getKafkaBrokerString(cfg.KafkaBrokers))

	return cfg
}

func getKafkaBrokerString(brokers []string) string {
	brokerStr := ""
	for i, broker := range brokers {
		brokerStr += broker
		if len(brokers) > i+1 {
			brokerStr += ","
		}
	}

	return brokerStr
}
