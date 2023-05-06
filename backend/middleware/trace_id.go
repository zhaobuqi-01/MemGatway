package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("TraceID", traceID)
		c.Next()
	}
}
