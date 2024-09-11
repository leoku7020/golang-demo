package logger

import (
	"context"
	"errors"
	"io/fs"
	"sync"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"

	"gs-mono/pkg/envkit"
)

const (
	mockEnv         = "Staging"
	mockPodName     = "example-api-main-6868d88fbd-bz8zv"
	mockProjName    = "example"
	mockServiceName = "example-api"
)

type loggerSuite struct {
	suite.Suite
}

func resetLogger() {
	regExporter = ExporterNone
	regLogger = nil
	registerOnce = sync.Once{}
	regDev = false
	regLevel = 0
}

func (s *loggerSuite) SetupSuite()    {}
func (s *loggerSuite) TearDownSuite() {}
func (s *loggerSuite) SetupTest() {
	envkit.ResetRegister()
	resetLogger()
}
func (s *loggerSuite) TearDownTest() {
	envkit.ResetRegister()
	resetLogger()
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(loggerSuite))
}

func (s *loggerSuite) TestRegister() {
	tests := []struct {
		Desc          string
		SetupTest     func(Exporter, bool)
		Exporter      Exporter
		Development   bool
		ExpRegistered bool
		ExpExporter   Exporter
		ExpFlush      error
	}{
		{
			Desc:          "not registered",
			ExpRegistered: false,
			ExpFlush:      ErrNotRegistered,
			ExpExporter:   ExporterNone,
		},
		{
			Desc: "empty implementation",
			SetupTest: func(e Exporter, b bool) {
				envkit.Register(envkit.Config{
					EnvNamespace: mockEnv,
					PodName:      mockPodName,
					ServiceName:  mockServiceName,
					ProjectName:  mockProjName,
				})
				Register(Config{
					Exporter:    e,
					Development: b,
				})
			},
			Exporter:      ExporterNone,
			ExpRegistered: true,
			ExpExporter:   ExporterNone,
			ExpFlush:      nil,
		},
		{
			Desc: "zap implementation",
			SetupTest: func(e Exporter, b bool) {
				envkit.Register(envkit.Config{
					EnvNamespace: mockEnv,
					PodName:      mockPodName,
					ServiceName:  mockServiceName,
					ProjectName:  mockProjName,
				})
				Register(Config{
					Exporter:    e,
					Development: b,
				})
			},
			Exporter:      ExporterZap,
			Development:   true,
			ExpRegistered: true,
			ExpExporter:   ExporterZap,
			// XXX: https://github.com/uber-go/zap/issues/991#issuecomment-962098428
			// It's the known issue mentioned in github. It didn't impact normal use case.
			// work around it in the future.
			ExpFlush: &fs.PathError{
				Op:   "sync",
				Path: "/dev/stderr",
				Err:  syscall.EBADF, // ENOTTY or EBADF
			},
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.SetupTest(t.Exporter, t.Development)
		}

		s.Require().Equal(t.ExpRegistered, Registered(), t.Desc)
		s.Require().Equal(t.ExpExporter, regExporter, t.Desc)
		s.Require().Equal(t.ExpFlush, Flush(), t.Desc)

		s.TearDownTest()
	}
}

func (s *loggerSuite) TestLog() {
	Register(Config{
		Exporter:    ExporterZap,
		Development: true,
	})

	Debug("debug")
	Info("info", WithFields(nil))
	Warn("warn", WithFields(Fields{"foo": "bar"}))
	Error("error", WithFields(Fields{"foo": "bar"}), WithError(errors.New("XD")))
}

func (s *loggerSuite) TestCtx() {
	Register(Config{
		Exporter:    ExporterZap,
		Development: true,
		CommonTags: map[string]string{
			"common": "tags",
		},
	})

	c := context.Background()
	c = ContextWithFields(c, Fields{"foo": "bar"})
	c = ContextWithFields(c, Fields{"number": 123})
	Ctx(c).Error("ctx error",
		WithFields(Fields{"bool": true, "foo": "bar2", "extra": "fields"}),
		WithError(errors.New("XD")),
		WithField("extra", "field"),
		WithField("append", 123),
	)
}

func (s *loggerSuite) TestPanic() {
	defer func() {
		r := recover()
		s.Require().NotNil(r)
	}()

	Register(Config{
		Exporter:    ExporterZap,
		Development: true,
	})

	Panic("panic")
}

func (s *loggerSuite) TestFatal() {
	s.T().Skip("it will trigger os.Exit(1), skip it")

	Register(Config{
		Exporter:    ExporterZap,
		Development: true,
	})

	Fatal("fatal")
}

func Example() {
	Register(Config{
		Exporter:    ExporterZap,
		Development: true,
	})

	// no output because it comes to stderr
	Debug("no stdout")

	// Output:
}
