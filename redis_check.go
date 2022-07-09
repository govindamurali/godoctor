package godoctor

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type redisChecker struct {
	client *redis.Client
}

const redisCheckerName = "redis"

func (c *redisChecker) Check(ctx context.Context, timeout time.Duration) error {
	errChan := make(chan error)
	go func() {
		_, err := c.client.Ping().Result()
		errChan <- err
	}()

	select {
	case <-time.After(timeout):
		return errors.New("ping timed out")
	case err := <-errChan:
		close(errChan)
		return err
	}
}

func (c *redisChecker) getName() checkerName {
	return redisCheckerName
}

func RedisChecker(client *redis.Client) IChecker {
	return &redisChecker{client: client}
}
