// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/pion/logging"
)

func main() {
	// Create a JSON logger factory
	factory := logging.NewJSONLoggerFactory(logging.WithJSONWriter(os.Stdout))

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
	// {"time":"2026-01-31T15:28:21.275282663-05:00","level":"INFO","msg":"API server started","scope":"api"}
	// {"time":"2026-01-31T15:28:21.275387414-05:00","level":"DEBUG","msg":"Processing request method=GET path=/users","scope":"api"}
	// {"time":"2026-01-31T15:28:21.275407888-05:00","level":"WARN","msg":"Rate limit approaching requests=95 limit=100","scope":"api"}
	// {"time":"2026-01-31T15:28:21.275426935-05:00","level":"INFO","msg":"Database connection established","scope":"database"}
	// {"time":"2026-01-31T15:28:21.275446032-05:00","level":"DEBUG","msg":"Executing query query=\"SELECT * FROM users\" duration_ms=15","scope":"database"}
	// {"time":"2026-01-31T15:28:21.27546512-05:00","level":"ERROR","msg":"Authentication failed user_id=12345 reason=invalid_token","scope":"auth"}
	// {"time":"2026-01-31T15:28:21.275483988-05:00","level":"INFO","msg":"User logged in user_id=67890 ip=192.168.1.100","scope":"auth"}
}
