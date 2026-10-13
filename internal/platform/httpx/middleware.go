package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)


// Logging gives every request an ID and logs are one JSON line when it finishes.

func Logging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = newID()
		}

		c.Header("X-Request-ID", reqID)
		c.Set("request_id", reqID)

		c.Next() 

		log.Info("request",
				"request_id", reqID,
				"method", c.Request.Method,
				"route", c.FullPath(),
				"status", c.Writer.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}


func newID() string {
	b  :=  make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
