package logger

import (
	"errors"
	"os"
	"strings"

	"github.com/airbrake/gobrake/v5"
	"go.uber.org/zap/zapcore"
)

var (
	repoName string
)

func init() {
	repoName = os.Getenv("REPO_NAME")
}

type Core struct {
	zapcore.LevelEnabler
	Notifier   *gobrake.Notifier
	coreFields map[string]interface{}
	depth      int
}

func NewZapBrakeCore(enab zapcore.LevelEnabler, notifier *gobrake.Notifier) (*Core, error) {
	if notifier == nil {
		return nil, errors.New("airbrake notifier is not defined")
	}
	core := &Core{
		LevelEnabler: enab,
		Notifier:     notifier,
		depth:        3,
	}
	return core, nil
}

// SetDepth method is for setting the depth of the notices
func (core *Core) SetDepth(depth int) {
	core.depth = depth
}

func (core *Core) With(fields []zapcore.Field) zapcore.Core {
	coreFields := make(map[string]interface{}, len(core.coreFields)+len(fields))
	for k, v := range core.coreFields {
		coreFields[k] = v
	}

	encoder := zapcore.NewMapObjectEncoder()
	// Process the fields passed directly.
	for _, field := range fields {
		field.AddTo(encoder)
	}

	for k, v := range encoder.Fields {
		coreFields[k] = v
	}

	core.coreFields = coreFields
	return core
}

func (core *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	parameters := make(map[string]interface{})
	notice := gobrake.NewNotice(entry.Message, nil, core.depth)
	for key, parameter := range core.coreFields {
		if key == "httpMethod" || key == "route" {
			notice.Context[key] = parameter
		} else {
			switch parameter := parameter.(type) {
			case error:
				parameters[key] = parameter.Error()
			default:
				parameters[key] = parameter
			}
		}
	}

	encoder := zapcore.NewMapObjectEncoder()
	// Process the fields passed directly.
	for _, field := range fields {
		field.AddTo(encoder)
	}

	for k, v := range encoder.Fields {
		parameters[k] = v
	}

	notice.Context["severity"] = entry.Level.String()
	notice.Params = parameters
	adjustProjectRootBasedOnRepoName(notice)
	core.Notifier.Notify(notice, nil)
	return nil
}

func (core *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if core.Enabled(entry.Level) {
		return checked.AddCore(entry, core)
	}
	return checked
}

func (core *Core) Sync() error {
	return nil
}

func adjustProjectRootBasedOnRepoName(notice *gobrake.Notice) {
	if repoName == "" {
		return
	}

	for i := 0; i < len(notice.Errors); i++ {
		err := notice.Errors[i]
		for j := 0; j < len(err.Backtrace); j++ {
			path := err.Backtrace[j].File
			if idx := strings.LastIndex(path, repoName); idx != -1 {
				err.Backtrace[j].File = path[idx+len(repoName):]
			}
		}
	}
}
