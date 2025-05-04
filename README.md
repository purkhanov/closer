# closer

🧹 Graceful shutdown helper for Go services.

This package provides a singleton-based mechanism to register and gracefully execute shutdown functions in order, especially useful for closing DB connections, stopping background workers, etc.

## 📦 Installation

```bash
go get github.com/purkhanov/closer
