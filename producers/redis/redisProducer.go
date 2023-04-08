package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type RedisProducer struct {
	client redis.Client
	Ttl    time.Duration
}

func (p *RedisProducer) Initialize(configFile string) {
	var options redis.Options

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load Redis configFile: %s", err)
	}

	err = json.Unmarshal(data, &options)
	if err != nil {
		log.Fatalf("Failed to parsa configuration parameters: %s", err)
	}

	p.client = *redis.NewClient(&options)
}

func (p *RedisProducer) Close() {
	err := p.client.Close()
	if err != nil {
		log.Fatalf("Failed to close Redis connection:\n%s", err)
	}
}

func (p *RedisProducer) Produce(k []byte, v []byte) {
	ctx := context.Background()
	err := p.client.Set(ctx, string(k), string(v), p.Ttl).Err()
	if err != nil {
		log.Fatalf("Failed to write data in Redis:\n%s", err)
	}
}
