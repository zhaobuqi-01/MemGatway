package mq

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type MessageQueue struct {
	redisClient       *redis.Client
	ctx               context.Context
	messageSeen       map[string]struct{}
	messageExpiration time.Duration
	counterMutex      sync.Mutex
}

func New(client *redis.Client, messageExpiration time.Duration) *MessageQueue {
	return &MessageQueue{
		redisClient:       client,
		ctx:               context.Background(),
		messageSeen:       make(map[string]struct{}),
		messageExpiration: messageExpiration,
		counterMutex:      sync.Mutex{},
	}
}

const defaultExpiration = 1 * time.Minute

func Default(client *redis.Client) *MessageQueue {
	return New(client, defaultExpiration)
}

// Publish 发布消息到指定频道
func (mq *MessageQueue) Publish(channel string, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = mq.redisClient.Publish(mq.ctx, channel, msg).Err()
	if err != nil {
		return err
	}

	return nil
}

func (mq *MessageQueue) Subscribe(channel string, deduplicate bool, callback func(channel string, message []byte)) error {
	pubsub := mq.redisClient.Subscribe(mq.ctx, channel)
	_, err := pubsub.Receive(mq.ctx)
	if err != nil {
		return err
	}

	channelMessage := pubsub.Channel()

	if deduplicate {
		callback = mq.deduplicate(callback)
	}

	go func() {
		for msg := range channelMessage {
			callback(msg.Channel, []byte(msg.Payload))
		}
	}()

	return nil
}

// Unsubscribe 取消订阅指定频道
func (mq *MessageQueue) Unsubscribe(channel string) error {
	pubsub := mq.redisClient.Subscribe(mq.ctx, channel)
	err := pubsub.Unsubscribe(mq.ctx, channel)
	if err != nil {
		return err
	}

	return nil
}

func (mq *MessageQueue) deduplicate(callback func(channel string, message []byte)) func(channel string, message []byte) {
	return func(channel string, message []byte) {
		messageStr := string(message)

		mq.counterMutex.Lock()
		_, seen := mq.messageSeen[messageStr]
		if !seen {
			mq.messageSeen[messageStr] = struct{}{}
			mq.counterMutex.Unlock()

			callback(channel, message)

			time.AfterFunc(mq.messageExpiration, func() {
				mq.counterMutex.Lock()
				delete(mq.messageSeen, messageStr)
				mq.counterMutex.Unlock()
			})
		} else {
			mq.counterMutex.Unlock()
		}
	}
}
