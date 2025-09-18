package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		fmt.Printf("[%s] %d %s %s - %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			method,
			path,
			duration,
		)
	}
}
