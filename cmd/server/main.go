package main

import (
	"net/http"
    // "log"
    // "time"
    // "fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/utrack/gin-csrf"

	// These paths must match your go.mod module name 🔗
	"garment-management-backend/internal/auth"
	"garment-management-backend/internal/models"

)

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().
// 			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		c.Writer.Header().
// 			Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		if c.Request.Method == "OPTIONS" {
// 			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}
// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082") // Best practice: Specify the origin 🎯
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-CSRF-TOKEN")		
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			// Change: Use 204 (No Content) or 200, and Abort to stop Gin from looking for a route 🛑
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

    store := cookie.NewStore([]byte("secret-key-for-session"))

    r.Use(sessions.Sessions("mysession", store))
    r.Use(CORSMiddleware())

    r.Use(csrf.Middleware(csrf.Options{
        Secret: "secret-key-for-csrf",
        ErrorFunc: func(c *gin.Context) {
            c.String(403, "CSRF token mismatch")
            c.Abort()
        },
    }))

	// Public Routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

    r.GET("/csrf-token", func(c *gin.Context) {
        c.JSON(200, gin.H{"token": csrf.GetToken(c)})
    })

	r.POST("/login", func(c *gin.Context) {
        var req models.LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        if auth.CheckCredentials(req.Username, req.Password) {
            // Returns both Access and Refresh tokens
            tokens, err := auth.GenerateTokenPair(req.Username)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
                return
            }
            c.JSON(http.StatusOK, tokens)
            return
        }
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    })

	r.POST("/refresh", func(c *gin.Context) {
        var body struct {
            RefreshToken string `json:"refresh_token" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
            return
        }

        // Validate the provided refresh token
        token, err := auth.ValidateToken(body.RefreshToken)
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        username := claims["username"].(string)

        // Issue a brand new pair 🎫
        newTokens, err := auth.GenerateTokenPair(username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not renew tokens"})
            return
        }

        c.JSON(http.StatusOK, newTokens)
    })

	// Protected Routes 🔒
	secure := r.Group("/secure")
	secure.Use(auth.JwtMiddleware())
	{
		secure.GET("", func(c *gin.Context) {
			user := c.GetString("username")
			c.JSON(http.StatusOK, gin.H{"message": "Welcome!", "user": user})
		})

        secure.POST("/logout", func(c *gin.Context) {
            // In a stateless JWT setup, we simply return success.
            // If you implement a "Blacklist" in Redis later, you would add the token to it here.
            c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
        })


	}

	r.Run(":8080")
}