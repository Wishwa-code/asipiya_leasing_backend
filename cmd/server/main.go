package main

import (
	"log"
	"net/http"

	// "time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	// These paths must match your go.mod module name 🔗
	"garment-management-backend/internal/auth"
	"garment-management-backend/internal/database"
	"garment-management-backend/internal/garmentOperation"
	"garment-management-backend/internal/leasing"
	"garment-management-backend/internal/middleware"
	"garment-management-backend/internal/models"
)

func main() {
	database.Connect()

	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}
	defer sqlDB.Close()

	r := gin.Default()
	middleware.StartLogWorker() // Start the background worker once 🪵
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AsyncOperationLogger()) // Use the async logger middleware ⚡

	// Serve uploaded files 📂
	r.Static("/uploads", "./uploads")

	r.NoRoute(middleware.CORSMiddleware(), func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route not found"})
	})

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
		user, success := auth.CheckCredentials(database.DB, req.Username, req.Password)

		if success {
			// Returns both Access and Refresh tokens
			tokens, err := auth.GenerateTokenPair(req.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
				return
			}

			fmt.Printf("🎉 Login Successful for user: %s\n", user.Name)

			c.JSON(http.StatusOK, gin.H{
				"tokens": tokens,
				"user":   user,
			})
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

	apiV1 := r.Group("/api")
	apiV1.Use(auth.JwtMiddleware())
	{
		// Automatically registers all /plantation routes under /api/v1 🌿
		garmentOperation.RegisterRoutes(apiV1)
		leasing.RegisterRoutes(apiV1)

		// Final URL will be: http://localhost:8080/api/v1/plantation/stats
	}

	r.Run(":8084")
}
