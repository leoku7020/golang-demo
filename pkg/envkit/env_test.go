package envkit

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	mockPodName     = "go-amazing-main-6868d88fbd-bz8zv"
	mockServiceName = "amazing-chatroom-rpc"
	mockProjectName = "amazing-chatroom"
)

type envSuite struct {
	suite.Suite
}

func (s *envSuite) SetupSuite()    {}
func (s *envSuite) TearDownSuite() {}
func (s *envSuite) SetupTest() {
	ResetRegister()
}
func (s *envSuite) TearDownTest() {
	ResetRegister()
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, new(envSuite))
}

func (s *envSuite) TestRegisterNoSuchEnv() {
	defer func() {
		r := recover()
		s.Require().NotNil(r)
		s.Require().Equal(errors.New("no such environment namespace"), r)
	}()

	Register(Config{EnvNamespace: "not existed"})
}

func (s *envSuite) TestRegisterNone() {
	defer func() {
		r := recover()
		s.Require().NotNil(r)
		s.Require().Equal(errors.New("no such environment namespace"), r)
	}()

	Register(Config{EnvNamespace: "None"})
}

func (s *envSuite) TestRegister() {
	tests := []struct {
		Desc            string
		SetupTest       func()
		ExpNamespace    Env
		ExpEnvNamespace string
		ExpPodName      string
		ExpServiceName  string
		ExpProjectName  string
	}{
		{
			Desc: "Normal register",
			SetupTest: func() {
				Register(Config{
					EnvNamespace: "production",
					PodName:      mockPodName,
					ServiceName:  mockServiceName,
					ProjectName:  mockProjectName,
				})
			},
			ExpEnvNamespace: EnvProduction.String(),
			ExpPodName:      mockPodName,
			ExpServiceName:  mockServiceName,
			ExpProjectName:  mockProjectName,
			ExpNamespace:    EnvProduction,
		},
		{
			Desc: "Register twice",
			SetupTest: func() {
				Register(Config{
					EnvNamespace: "production",
					PodName:      mockPodName,
					ServiceName:  mockServiceName,
					ProjectName:  mockProjectName,
				})
				Register(Config{
					EnvNamespace: "staging",
					PodName:      "another-pod-name",
					ProjectName:  "another-project-name",
					ServiceName:  "another-service-name",
				})
			},
			ExpEnvNamespace: EnvProduction.String(),
			ExpPodName:      mockPodName,
			ExpServiceName:  mockServiceName,
			ExpProjectName:  mockProjectName,
			ExpNamespace:    EnvProduction,
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.SetupTest()
		}

		s.Require().Equal(t.ExpNamespace, Namespace(), t.Desc)
		s.Require().Equal(t.ExpEnvNamespace, EnvNamespace(), t.Desc)
		s.Require().Equal(t.ExpPodName, PodName(), t.Desc)
		s.Require().Equal(t.ExpServiceName, ServiceName(), t.Desc)
		s.Require().Equal(t.ExpProjectName, ProjectName(), t.Desc)

		s.TearDownTest()
	}
}

func (s *envSuite) TestGetByEnvVariable() {
	tests := []struct {
		Desc            string
		SetupTest       func()
		ExpEnvNamespace string
		ExpPodName      string
		ExpServiceName  string
		ExpProjectName  string
	}{
		{
			Desc: "Normal case",
			SetupTest: func() {
				s.Require().NoError(os.Setenv("ENV_NAMESPACE", "development"))
				s.Require().NoError(os.Setenv("ENV_POD_NAME", mockPodName))
				s.Require().NoError(os.Setenv("ENV_SERVICE_NAME", mockServiceName))
				s.Require().NoError(os.Setenv("ENV_PROJECT_NAME", mockProjectName))
			},
			ExpEnvNamespace: EnvDevelopment.String(),
			ExpPodName:      mockPodName,
			ExpServiceName:  mockServiceName,
			ExpProjectName:  mockProjectName,
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.SetupTest()
		}

		s.Require().Equal(t.ExpEnvNamespace, EnvNamespace(), t.Desc)
		s.Require().Equal(t.ExpPodName, PodName(), t.Desc)
		s.Require().Equal(t.ExpServiceName, ServiceName(), t.Desc)
		s.Require().Equal(t.ExpProjectName, ProjectName(), t.Desc)

		s.TearDownTest()
	}
}

func (s *envSuite) TestRetrieveCallerInfo() {
	info := RetrieveCallerInfo()
	s.Require().Equal("env_test.go", info.FileName)
	s.Require().Equal("envkit", info.PkgName)
	s.Require().Equal("gs-mono/pkg/envkit", info.FullPkgName)
	s.Require().Equal("TestRetrieveCallerInfo", info.FuncName)
	s.Require().Equal("(*envSuite).TestRetrieveCallerInfo", info.FullFuncName)
	s.Require().Equal(158, info.Line) // it's the line of `info := RetrieveCallerInfo()`
}
