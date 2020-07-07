package config

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestNewKernel(t *testing.T) {
	cfg, err := NewKernel(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", cfg)

	// Verify kafka broker env variable
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		t.Error(errors.New("no kafka brokers found"))
	}
	t.Log(kafkaBrokers)
}
