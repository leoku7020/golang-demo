// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package ctxkit

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// KeyNone is a Key of type None.
	// Not existed
	KeyNone Key = iota
	// KeyLogger is a Key of type Logger.
	// Used in pkg/logger to store several key-value pairs in the context
	KeyLogger
)

var ErrInvalidKey = errors.New("not a valid Key")

const _KeyName = "NoneLogger"

var _KeyMap = map[Key]string{
	KeyNone:   _KeyName[0:4],
	KeyLogger: _KeyName[4:10],
}

// String implements the Stringer interface.
func (x Key) String() string {
	if str, ok := _KeyMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Key(%d)", x)
}

var _KeyValue = map[string]Key{
	_KeyName[0:4]:                   KeyNone,
	strings.ToLower(_KeyName[0:4]):  KeyNone,
	_KeyName[4:10]:                  KeyLogger,
	strings.ToLower(_KeyName[4:10]): KeyLogger,
}

// ParseKey attempts to convert a string to a Key.
func ParseKey(name string) (Key, error) {
	if x, ok := _KeyValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _KeyValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Key(0), fmt.Errorf("%s is %w", name, ErrInvalidKey)
}
