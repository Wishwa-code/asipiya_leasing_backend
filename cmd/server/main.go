package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	// These paths must match your go.mod module name 🔗
	"basic/internal/auth"
	"basic/internal/models"
)

func main() {
	r := gin.Default()

	// Public Routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
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
	}

	r.Run(":8080")
}