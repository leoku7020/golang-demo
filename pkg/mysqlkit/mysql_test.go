package mysqlkit

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type mysqlSuite struct {
	suite.Suite
}

func (s *mysqlSuite) SetupSuite() {}

func (s *mysqlSuite) TearDownSuite() {}

func (s *mysqlSuite) SetupTest() {}

func (s *mysqlSuite) TearDownTest() {}

func TestMySQLSuite(t *testing.T) {
	suite.Run(t, new(mysqlSuite))
}
