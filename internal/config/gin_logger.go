package config

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Now().UnixMilli() - start.UnixMilli()

		entry := logrus.WithFields(logrus.Fields{
			"duration_ms": duration,
			"method":      c.Request.Method,
			"path":        c.Request.RequestURI,
			"status":      c.Writer.Status(),
			"referrer":    c.Request.Referer(),
			"request_id":  c.Writer.Header().Get("Request-Id"),
		})

		if c.Writer.Status() >= http.StatusInternalServerError {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("request")
		}
	}
}
