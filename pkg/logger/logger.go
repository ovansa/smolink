package logger

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Path: %s | Duration: %v",
			ctx.Request.Method, ctx.Writer.Status(), ctx.Request.URL.Path, duration)
	}
}
