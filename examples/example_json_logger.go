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
	apiLogger.Debugf("Processing request method=%s path=%s", "GET", "/users")
	apiLogger.Warnf("Rate limit approaching requests=%d limit=%d", 95, 100)

	dbLogger.Info("Database connection established")
	dbLogger.Debugf("Executing query query=%q duration_ms=%d", "SELECT * FROM users", 15)

	authLogger.Errorf("Authentication failed user_id=%s reason=%s", "12345", "invalid_token")
	authLogger.Infof("User logged in user_id=%s ip=%s", "67890", "192.168.1.100")

	// nolint:lll
	// Example output will be JSON formatted like:
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"API server started","scope":"api"}
	// {"time":"2023-12-07T10:30:00Z","level":"DEBUG","msg":"Processing request","scope":"api","method":"GET","path":"/users"}
	// {"time":"2023-12-07T10:30:00Z","level":"WARN","msg":"Rate limit approaching","scope":"api","requests":95,"limit":100}
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"Database connection established","scope":"database"}
	// {"time":"2023-12-07T10:30:00Z","level":"DEBUG","msg":"Executing query","scope":"database","query":"SELECT * FROM users","duration_ms":15}
	// {"time":"2023-12-07T10:30:00Z","level":"ERROR","msg":"Authentication failed","scope":"auth","user_id":"12345","reason":"invalid_token"}
	// {"time":"2023-12-07T10:30:00Z","level":"INFO","msg":"User logged in","scope":"auth","user_id":"67890","ip":"192.168.1.100"}
}
