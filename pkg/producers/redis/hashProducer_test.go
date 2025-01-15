//go:build exclude

// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHashProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer := &HashProducer{}
	producer.Initialize(configFile)

	assert.NotNil(t, producer.client, "Redis client should be initialized")
}

func TestHashProducer_Produce(t *testing.T) {
	configFile := "config.json.example"

	producer := &HashProducer{
		Ttl: time.Minute,
	}
	producer.Initialize(configFile)

	ctx := context.Background()
	key := "test_hash_key"
	value := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	valueBytes, _ := json.Marshal(value)

	producer.Produce(ctx, []byte(key), valueBytes, nil)

	// Verify the data in Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust as necessary
	})
	defer client.Close()

	result, err := client.HGetAll(ctx, key).Result()
	assert.NoError(t, err, "Should not error when getting hash from Redis")
	assert.Equal(t, "value1", result["field1"], "Field1 should match")
	assert.Equal(t, "value2", result["field2"], "Field2 should match")
}

func TestHashProducer_Close(t *testing.T) {
	configFile := "config.json.example"

	producer := &HashProducer{}
	producer.Initialize(configFile)

	err := producer.Close(context.Background())
	assert.NoError(t, err, "Should not error when closing Redis connection")
}
