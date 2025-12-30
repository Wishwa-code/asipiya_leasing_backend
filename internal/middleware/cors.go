package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing settings 🌐
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use c.Header instead of c.Writer.Header().Set to ensure persistence during Aborts 🚀
		c.Header("Access-Control-Allow-Origin", "http://localhost:8082")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-XSRF-TOKEN, X-CSRF-TOKEN")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

