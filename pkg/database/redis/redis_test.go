package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectRedis(t *testing.T) {
	// Connect to the Redis database
	client, err := connectRedis()
	if err != nil {
		t.Errorf("Error connecting to Redis database: %v", err)
	}

	// Make sure the connection is not nil
	assert.NotNil(t, client)

	// Make sure we can set and retrieve a key-value pair
	key := "test_key"
	value := "test_value"
	err = client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		t.Errorf("Error setting key-value pair in Redis: %v", err)
	}
	result, err := client.Get(context.Background(), key).Result()
	if err != nil {
		t.Errorf("Error retrieving key-value pair from Redis: %v", err)
	}
	assert.Equal(t, value, result)

	// Close the Redis connection
	err = client.Close()
	if err != nil {
		t.Errorf("Error closing Redis connection: %v", err)
	}
}
