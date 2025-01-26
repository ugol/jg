//go:build exclude

package wamprpc

import (
	"context"
	"testing"

	"github.com/jrnd-io/jr/pkg/producers/wamprpc"
)

func TestProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer, err := wamprpc.ProducerFactory("wamprpc")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
}

func TestProducer_Close(t *testing.T) {
	configFile := "config.json.example"

	producer, err := wamprpc.ProducerFactory("wamprpc")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	producer.Close()
}

func TestProducer_Produce(t *testing.T) {
	configFile := "config.json.example"

	producer, err := wamprpc.ProducerFactory("wamprpc")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error initializing producer: %v", err)
	}

	ctx := context.Background()
	key := "loo"
	val := "foo"
	exp := 0
	producer.Produce(ctx, key, val, exp)

}
