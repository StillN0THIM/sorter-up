package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Printf("[%d] %s %s - %s",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			time.Since(start))
	}
}
