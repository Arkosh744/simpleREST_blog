package rest

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// to process the request
		c.Next()
		// End time
		endTime := time.Now()
		// Execution time
		latencyTime := endTime.Sub(startTime)
		// Request method
		reqMethod := c.Request.Method
		// Request routing
		reqUri := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// The requestIP
		clientIP := c.ClientIP()
		// Log format
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode, latencyTime, clientIP, reqMethod, reqUri)
	}
}
