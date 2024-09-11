//go:generate go-enum -f=$GOFILE --nocase

package ctxkit

// Key is an enumeration of keys used to specify .
/*
ENUM(
None // Not existed
Logger // Used in pkg/logger to store several key-value pairs in the context
)
*/
type Key int32
