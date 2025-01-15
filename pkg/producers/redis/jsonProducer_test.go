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

func TestJSONProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer := &JSONProducer{}
	producer.Initialize(configFile)

	assert.NotNil(t, producer.client, "Redis client should be initialized")
}

func TestJSONProducer_Produce(t *testing.T) {
	configFile := "config.json.example"

	producer := &JSONProducer{
		Ttl: time.Minute,
	}
	producer.Initialize(configFile)

	ctx := context.Background()
	key := "test_json_key"
	// Create a test JSON object with nested structures
	testJSON := `{
        "id": "2210",
        "user": {
            "name": "Foogaro",
            "year": 1978,
            "email": "luigi@foogaro.com"
        }
    }`
	producer.Produce(ctx, []byte(key), []byte(testJSON), nil)

	// Verify the data in Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust as necessary
	})
	defer client.Close()

	// Use JSON.GET to retrieve the stored JSON
	result, err := client.Do(ctx, "JSON.GET", key, "$").Text()
	assert.NoError(t, err, "Should not error when getting JSON from Redis")

	// Compare the JSON strings (after normalizing them)
	var expected, actual interface{}
	err = json.Unmarshal([]byte(testJSON), &expected)
	assert.NoError(t, err, "Should parse expected JSON")

	err = json.Unmarshal([]byte(result), &actual)
	assert.NoError(t, err, "Should parse actual JSON")

	assert.Equal(t, expected, actual, "Stored JSON should match original")
}

func TestJSONProducer_Close(t *testing.T) {
	configFile := "config.json.example"

	producer := &JSONProducer{}
	producer.Initialize(configFile)

	err := producer.Close(context.Background())
	assert.NoError(t, err, "Should not error when closing Redis connection")
}
