package middleware

import (
	"log"
	"smolink/internal/errors"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Path: %s | Duration: %v",
			ctx.Request.Method, ctx.Writer.Status(), ctx.Request.URL.Path, duration)
	}
}

func respondWithError(ctx *gin.Context, err error) {
	apiErr := errors.ExtractAPIError(err)
	ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Handle panics
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				respondWithError(ctx, errors.ErrInternal)
			}
		}()

		ctx.Next()

		// Process any errors added to the context
		if len(ctx.Errors) > 0 {
			for _, ginErr := range ctx.Errors {
				log.Printf("Error: %v", ginErr.Err)
			}

			// Get the last error
			lastErr := ctx.Errors.Last()
			respondWithError(ctx, lastErr.Err)
		}
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
