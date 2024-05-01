package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// gin ratelimit example
// https://github.com/gin-gonic/examples/blob/master/ratelimiter/rate.go
func leakyBucket(rps int) gin.HandlerFunc {
	limit := ratelimit.New(rps)
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		prev = now
		_ = prev
	}
}
