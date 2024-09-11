// Package bootkit plays as a role of manager coordinating booting up or shutting down the system.
// It provides convenient interface for working with os.Signal.
// Multiple hooks with level can be applied, they will be called simultaneously with the same level on app shutdown.
package bootkit
