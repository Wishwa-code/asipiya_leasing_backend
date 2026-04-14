package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("my_super_secret_key")

// TokenPair holds both short-lived and long-lived tokens 🎫
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair creates an Access Token (15m) and a Refresh Token (7d)
func GenerateTokenPair(username string) (TokenPair, error) {
	// 1. Access Token - Short expiration for security ⏱️
	accessToken, err := createToken(username, time.Minute*15)
	if err != nil {
		return TokenPair{}, err
	}

	// 2. Refresh Token - Long expiration to keep user logged in 🔄
	refreshToken, err := createToken(username, time.Hour*24*7)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Helper function to minimize code duplication when creating tokens
func createToken(username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtSecret)
}

// ValidateToken parses a token string and returns the token object if valid
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
}

// JwtMiddleware handles the 🛡️ protection of routes
// func JwtMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		// Use the new ValidateToken helper to keep logic DRY 🌵
// 		token, err := ValidateToken(tokenString)

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		c.Set("username", claims["username"])
// 		c.Next()
// 	}
// }

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Log the incoming request path and if header exists 📝
		log.Printf("[JWT] Checking auth for: %s %s", c.Request.Method, c.Request.URL.Path)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("[JWT] ❌ Missing or malformed Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ValidateToken(tokenString)

		if err != nil {
			log.Printf("[JWT] ❌ Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Printf("[JWT] ❌ Token is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		log.Printf("[JWT] ✅ Authorized user: %v", claims["username"])

		c.Set("username", claims["username"])
		c.Next()
	}
}
