// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/pion/logging"
)

func main() {
	// Create a JSON logger factory
	factory := logging.NewJSONLoggerFactory()
	factory.Writer = os.Stdout // Output to stdout for this example

	// Create loggers for different scopes
	apiLogger := factory.NewLogger("api")
	dbLogger := factory.NewLogger("database")
	authLogger := factory.NewLogger("auth")

	// Log some messages
	apiLogger.Info("API server started")
	apiLogger.Debug("Processing request", "method", "GET", "path", "/users")
	apiLogger.Warn("Rate limit approaching", "requests", 95, "limit", 100)

	dbLogger.Info("Database connection established")
	dbLogger.Debug("Executing query", "query", "SELECT * FROM users", "duration_ms", 15)

	authLogger.Error("Authentication failed", "user_id", "12345", "reason", "invalid_token")
	authLogger.Info("User logged in", "user_id", "67890", "ip", "192.168.1.100")

	// Example output will be JSON formatted like:
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"API server started","scope":"api"}
	// {"time":"2023-12-07T10:30:00Z","level":"DEBUG","msg":"Processing request","scope":"api","method":"GET","path":"/users"}
	// {"time":"2023-12-07T10:30:00Z","level":"WARN","msg":"Rate limit approaching","scope":"api","requests":95,"limit":100}
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"Database connection established","scope":"database"}
	// {"time":"2023-12-07T10:30:00Z","level":"DEBUG","msg":"Executing query","scope":"database","query":"SELECT * FROM users","duration_ms":15}
	// {"time":"2023-12-07T10:30:00Z","level":"ERROR","msg":"Authentication failed","scope":"auth","user_id":"12345","reason":"invalid_token"}
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"User logged in","scope":"auth","user_id":"67890","ip":"192.168.1.100"}
} 