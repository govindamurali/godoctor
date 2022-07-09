package godoctor

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type IChecker interface {
	Check(ctx context.Context, timeout time.Duration) error
	getName() checkerName
}

type checkerName string

func GetHandler(timeout time.Duration, checkers ...IChecker) gin.HandlerFunc {
	return func(context *gin.Context) {
		success := true
		response := map[string]interface{}{}
		wg := sync.WaitGroup{}
		mutex := sync.RWMutex{}
		wg.Add(len(checkers))
		for index := range checkers {
			checker := checkers[index]
			go func() {
				defer wg.Done()
				err := checker.Check(context, timeout)
				mutex.Lock()
				defer mutex.Unlock()
				if err != nil {
					success = false
					response[string(checker.getName())] = err.Error()
				} else {
					response[string(checker.getName())] = "OK"
				}
			}()
		}
		wg.Wait()
		if success {
			context.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"data":    response,
			})
		} else {
			context.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"success": false,
				"error":   response,
			})
		}

	}

}
