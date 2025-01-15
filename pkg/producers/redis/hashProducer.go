package redis

import (
    "context"
    "encoding/json"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog/log"
)

type HashProducer struct {
    client redis.Client
    Ttl    time.Duration
}

func (p *HashProducer) Initialize(configFile string) {
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

func (p *HashProducer) Close(_ context.Context) error {
    err := p.client.Close()
    if err != nil {
        log.Warn().Err(err).Msg("Failed to close Redis connection")
    }
    return err
}

func (p *HashProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
    // Parse the JSON value into a map
    var fields map[string]interface{}
    err := json.Unmarshal(v, &fields)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to unmarshal JSON into hash fields")
    }

    // Use HSet to set multiple hash fields at once
    err = p.client.HSet(ctx, string(k), fields).Err()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to write hash data to Redis")
    }

    // Set TTL if specified
    if p.Ttl > 0 {
        err = p.client.Expire(ctx, string(k), p.Ttl).Err()
        if err != nil {
            log.Fatal().Err(err).Msg("Failed to set TTL on Redis hash")
        }
    }
} 