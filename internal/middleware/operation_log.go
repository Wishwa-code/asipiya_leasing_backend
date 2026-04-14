package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationLog represents the data we want to save
type OperationLog struct {
	Username  string        `json:"username"`
	Method    string        `json:"method"`
	Path      string        `json:"path"`
	Status    int           `json:"status"`
	Latency   time.Duration `json:"latency"`
	IP        string        `json:"ip"`
	Timestamp time.Time     `json:"timestamp"`
}

// logQueue is a buffered channel that acts as a waiting room for logs.
// Size 10000 allows a large burst of logs before the app starts blocking.
var logQueue = make(chan OperationLog, 10000)

// StartLogWorker runs in the background and processes logs from the queue.
// This should be called once when the server starts.
func StartLogWorker() {
	go func() {
		for logEntry := range logQueue {
			// Process the log (Save to MySQL/Postgres or File)
			// Using slog for structured output to console/file
			slog.Info("audit_log",
				"username", logEntry.Username,
				"method", logEntry.Method,
				"path", logEntry.Path,
				"status", logEntry.Status,
				"latency", logEntry.Latency,
				"ip", logEntry.IP,
				"timestamp", logEntry.Timestamp,
			)

			// If saving to DB, do it here.
			// db.Create(&logEntry)
		}
	}()
}

// AsyncOperationLogger is a middleware that captures request data and sends it to the log queue asynchronously.
func AsyncOperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next() // Wait for the rest of the request to finish

		// Capture data AFTER request is done
		// Note: We use "username" because that's what JwtMiddleware sets.
		username := c.GetString("username")
		if username == "" {
			username = "anonymous"
		}

		entry := OperationLog{
			Username:  username,
			Method:    c.Request.Method,
			Path:      path,
			Status:    c.Writer.Status(),
			Latency:   time.Since(start),
			IP:        c.ClientIP(),
			Timestamp: time.Now(),
		}

		// Blocking send to the channel to ensure every log is recorded 🛡️
		logQueue <- entry
	}
}
