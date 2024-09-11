//go:generate go-enum -f=$GOFILE --nocase

package logger

// Level is an enumeration of levels.
/*
ENUM(
Debug=-1 // Debug logs are typically voluminous, and are usually disabled in production.
Info // Info is the default logging priority.
Warn // Warn logs are more important than Info, but don't need individual human review.
Error // Error logs are high-priority. If an application is running smoothly. it shouldn't generate any error-level logs.
DPanic // DPanic logs are particularly important errors. In development the logger panics after writing the message.
Panic // Panic logs a message, then panics.
Fatal // Fatal logs a message, then calls os.Exit(1).
)
*/
type Level int8

func (l *Level) MarshalFlag() (string, error) {
	return l.String(), nil
}

func (l *Level) UnmarshalFlag(value string) error {
	lvl, err := ParseLevel(value)
	if err != nil {
		return err
	}

	*l = lvl
	return nil
}
