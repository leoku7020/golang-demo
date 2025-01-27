// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package logger

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// LevelDebug is a Level of type Debug.
	// Debug logs are typically voluminous, and are usually disabled in production.
	LevelDebug Level = iota + -1
	// LevelInfo is a Level of type Info.
	// Info is the default logging priority.
	LevelInfo
	// LevelWarn is a Level of type Warn.
	// Warn logs are more important than Info, but don't need individual human review.
	LevelWarn
	// LevelError is a Level of type Error.
	// Error logs are high-priority. If an application is running smoothly. it shouldn't generate any error-level logs.
	LevelError
	// LevelDPanic is a Level of type DPanic.
	// DPanic logs are particularly important errors. In development the logger panics after writing the message.
	LevelDPanic
	// LevelPanic is a Level of type Panic.
	// Panic logs a message, then panics.
	LevelPanic
	// LevelFatal is a Level of type Fatal.
	// Fatal logs a message, then calls os.Exit(1).
	LevelFatal
)

var ErrInvalidLevel = errors.New("not a valid Level")

const _LevelName = "DebugInfoWarnErrorDPanicPanicFatal"

var _LevelMap = map[Level]string{
	LevelDebug:  _LevelName[0:5],
	LevelInfo:   _LevelName[5:9],
	LevelWarn:   _LevelName[9:13],
	LevelError:  _LevelName[13:18],
	LevelDPanic: _LevelName[18:24],
	LevelPanic:  _LevelName[24:29],
	LevelFatal:  _LevelName[29:34],
}

// String implements the Stringer interface.
func (x Level) String() string {
	if str, ok := _LevelMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Level(%d)", x)
}

var _LevelValue = map[string]Level{
	_LevelName[0:5]:                    LevelDebug,
	strings.ToLower(_LevelName[0:5]):   LevelDebug,
	_LevelName[5:9]:                    LevelInfo,
	strings.ToLower(_LevelName[5:9]):   LevelInfo,
	_LevelName[9:13]:                   LevelWarn,
	strings.ToLower(_LevelName[9:13]):  LevelWarn,
	_LevelName[13:18]:                  LevelError,
	strings.ToLower(_LevelName[13:18]): LevelError,
	_LevelName[18:24]:                  LevelDPanic,
	strings.ToLower(_LevelName[18:24]): LevelDPanic,
	_LevelName[24:29]:                  LevelPanic,
	strings.ToLower(_LevelName[24:29]): LevelPanic,
	_LevelName[29:34]:                  LevelFatal,
	strings.ToLower(_LevelName[29:34]): LevelFatal,
}

// ParseLevel attempts to convert a string to a Level.
func ParseLevel(name string) (Level, error) {
	if x, ok := _LevelValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _LevelValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Level(0), fmt.Errorf("%s is %w", name, ErrInvalidLevel)
}
