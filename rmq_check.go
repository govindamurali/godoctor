package godoctor

import (
	"context"
	"errors"
	"github.com/streadway/amqp"
	"time"
)

type rmqChecker struct {
	userName string
	password string
	host     string
	port     string
}

const rabbitMqCheckerName = "rabbit_mq"

func (c *rmqChecker) getName() checkerName {
	return rabbitMqCheckerName
}

func (c *rmqChecker) Check(ctx context.Context, timeout time.Duration) error {
	errChan := make(chan error)
	go func() {
		conn, err := amqp.DialConfig("amqp://"+c.userName+":"+c.password+"@"+c.host+":"+c.port+"/", amqp.Config{
			Heartbeat: 2 * timeout,
			Locale:    "en_US",
		})
		if err != nil {
			errChan <- err
		}
		errChan <- conn.Close()
	}()
	select {
	case <-time.After(timeout):
		return errors.New("ping timed out")
	case err := <-errChan:
		close(errChan)
		return err
	}
}

func RabbitMqChecker(userName, password, host, port string) IChecker {
	return &rmqChecker{
		userName: userName,
		password: password,
		host:     host,
		port:     port,
	}
}
