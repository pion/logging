<h1 align="center">
  <br>
  Pion Logging
  <br>
</h1>
<h4 align="center">The Pion logging library</h4>
<p align="center">
  <a href="https://pion.ly"><img src="https://img.shields.io/badge/pion-logging-gray.svg?longCache=true&colorB=brightgreen" alt="Pion transport"></a>
  <a href="https://discord.gg/PngbdqpFbt"><img src="https://img.shields.io/badge/join-us%20on%20discord-gray.svg?longCache=true&logo=discord&colorB=brightblue" alt="join us on Discord"></a> <a href="https://bsky.app/profile/pion.ly"><img src="https://img.shields.io/badge/follow-us%20on%20bluesky-gray.svg?longCache=true&logo=bluesky&colorB=brightblue" alt="Follow us on Bluesky"></a> 
  <br>
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/pion/logging/test.yaml">
  <a href="https://pkg.go.dev/github.com/pion/logging"><img src="https://pkg.go.dev/badge/github.com/pion/logging.svg" alt="Go Reference"></a>
  <a href="https://codecov.io/gh/pion/logging"><img src="https://codecov.io/gh/pion/logging/branch/master/graph/badge.svg" alt="Coverage Status"></a>
  <a href="https://goreportcard.com/report/github.com/pion/logging"><img src="https://goreportcard.com/badge/github.com/pion/logging" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>
<br>

## Features

- **Text Logging**: Traditional text-based logging with customizable prefixes and levels
- **JSON Logging**: Structured JSON logging using Go's `slog` library (Go 1.21+)
- **Level-based Filtering**: Support for TRACE, DEBUG, INFO, WARN, ERROR, and DISABLED levels
- **Scope-based Configuration**: Different log levels for different scopes/components
- **Environment Variable Configuration**: Configure log levels via environment variables
- **Thread-safe**: All logging operations are thread-safe

## Usage

### Text Logging (Default)

```go
import "github.com/pion/logging"

// Create a logger factory
factory := logging.NewDefaultLoggerFactory()

// Create loggers for different scopes
apiLogger := factory.NewLogger("api")
dbLogger := factory.NewLogger("database")

// Log messages
apiLogger.Info("API server started")
apiLogger.Debug("Processing request")
dbLogger.Error("Database connection failed")
```

### JSON Logging

```go
import "github.com/pion/logging"

// Create a JSON logger factory
factory := logging.NewJSONLoggerFactory()

// Create loggers for different scopes
apiLogger := factory.NewLogger("api")
dbLogger := factory.NewLogger("database")

// Log messages with structured data
apiLogger.Info("API server started")
apiLogger.Debug("Processing request", "method", "GET", "path", "/users")
dbLogger.Error("Database connection failed", "error", "connection timeout")
```

### Environment Variable Configuration

Set environment variables to configure log levels:

```bash
# Enable all log levels
export PION_LOG_TRACE=all
export PION_LOG_DEBUG=all
export PION_LOG_INFO=all
export PION_LOG_WARN=all
export PION_LOG_ERROR=all

# Enable specific scopes
export PION_LOG_DEBUG=api,database
export PION_LOG_INFO=feature1,feature2
```

### Roadmap
The library is used as a part of our WebRTC implementation. Please refer to that [roadmap](https://github.com/pion/webrtc/issues/9) to track our major milestones.

### Community
Pion has an active community on the [Discord](https://discord.gg/PngbdqpFbt).

Follow the [Pion Bluesky](https://bsky.app/profile/pion.ly) or [Pion Twitter](https://twitter.com/_pion) for project updates and important WebRTC news.

We are always looking to support **your projects**. Please reach out if you have something to build!
If you need commercial support or don't want to use public methods you can contact us at [team@pion.ly](mailto:team@pion.ly)

### Contributing
Check out the [contributing wiki](https://github.com/pion/webrtc/wiki/Contributing) to join the group of amazing people making this project possible

### License
MIT License - see [LICENSE](LICENSE) for full text
