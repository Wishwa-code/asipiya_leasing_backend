package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing settings 🌐
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use c.Header instead of c.Writer.Header().Set to ensure persistence during Aborts 🚀
		allowedOrigins := map[string]bool{
			"http://localhost:5175":               true,
			"https://vite-vue-two-rho.vercel.app": true,
		}

		// Check if the current origin is allowed
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
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

// package middleware

// import (
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// // CORSMiddleware handles Cross-Origin Resource Sharing settings using the built-in package 🌐
// func CORSMiddleware() gin.HandlerFunc {
// 	config := cors.DefaultConfig()

// 	// Add your specific origins 🔒
// 	config.AllowOrigins = []string{
// 		"http://localhost:8082",
// 		"https://vite-vue-two-rho.vercel.app",
// 	}

// 	// Set allowed methods and headers 🛠️
// 	config.AllowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"}
// 	config.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "X-XSRF-TOKEN", "X-CSRF-TOKEN"}
// 	config.AllowCredentials = true

// 	return cors.New(config)
// }
