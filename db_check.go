package godoctor

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type dbChecker struct {
	db *sql.DB
}

const postgres = "database"

func (c *dbChecker) getName() checkerName {
	return postgres
}

func (c *dbChecker) Check(ctx context.Context, timeout time.Duration) error {
	errChan := make(chan error)
	go func() {
		errChan <- c.db.PingContext(ctx)
	}()
	select {
	case <-time.After(timeout):
		return errors.New("ping timed out")
	case err := <-errChan:
		close(errChan)
		return err
	}
}

func DbChecker(db *sql.DB) IChecker {
	return &dbChecker{db: db}
}
