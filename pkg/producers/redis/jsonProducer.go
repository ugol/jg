package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type JSONProducer struct {
	client redis.Client
	Ttl    time.Duration
}

func (p *JSONProducer) Initialize(configFile string) {
	var options redis.Options

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load Redis configFile")
	}

	err = json.Unmarshal(data, &options)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	p.client = *redis.NewClient(&options)
}

func (p *JSONProducer) Close(_ context.Context) error {
	err := p.client.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close Redis connection")
	}
	return err
}

func (p *JSONProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	// Verify the input is valid JSON
	var jsonData interface{}
	err := json.Unmarshal(v, &jsonData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to validate JSON data")
	}

	// Use JSON.SET to store the JSON document
	err = p.client.Do(ctx, "JSON.SET", string(k), "$", string(v)).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write JSON data to Redis")
	}

	// Set TTL if specified
	if p.Ttl > 0 {
		err = p.client.Expire(ctx, string(k), p.Ttl).Err()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to set TTL on Redis key")
		}
	}
}
