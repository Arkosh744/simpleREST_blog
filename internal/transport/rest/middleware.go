package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Log format
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode, latencyTime, clientIP, reqMethod, reqUri)
	}
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromContex(c)
		if err != nil {
			log.Println("authMiddleware", err)
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userId, err := h.usersService.ParseToken(c, token)
		if err != nil {
			log.Println("authMiddleware", err)
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		// Set context value
		c.Set(string(rune(ctxUserID)), userId)
		c.Next()

	}
}

func getTokenFromContex(c *gin.Context) (string, error) {
	header := c.Request.Header["Authorization"]
	if len(header) == 0 {
		return "", errors.New("no token in header")
	}
	log.Println("header", header)
	if header[0] == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header[0], " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
