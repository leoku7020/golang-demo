package sql

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"gs-mono/pkg/dockerkit"
	"gs-mono/pkg/logger"
	"gs-mono/pkg/mysqlkit"
)

const (
	migrationDir = "../../../database/migrations/example"
	sqlURLFormat = "root:root@tcp(localhost:%s)/mysql?charset=utf8mb4&parseTime=True"
)

var (
	mockCTX = context.Background()
)

type gooseSuite struct {
	suite.Suite

	db        *gorm.DB
	mysqlPort string
}

func (s *gooseSuite) migrationDir() string {
	return migrationDir
}

func (s *gooseSuite) mysqlURL() string {
	return fmt.Sprintf(sqlURLFormat, s.mysqlPort)
}

func (s *gooseSuite) SetupSuite() {
	// setup logger
	logger.Register(logger.Config{
		Exporter:    logger.ExporterZap,
		Development: true,
	})

	ports, err := dockerkit.RunExtDockers(mockCTX, []dockerkit.Image{
		dockerkit.ImageMySQL,
	})
	s.Require().NoError(err)
	s.mysqlPort = ports[0]

	db, err := mysqlkit.NewGORM(
		s.mysqlURL(),
	)
	s.Require().NoError(err)
	s.db = db
}

func (s *gooseSuite) TearDownSuite() {
	dockerkit.PurgeExtDockers(mockCTX, []dockerkit.Image{
		dockerkit.ImageMySQL,
	})

	logger.Flush()
}

func (s *gooseSuite) SetupTest() {}

func (s *gooseSuite) TearDownTest() {}

func TestGooseSuite(t *testing.T) {
	suite.Run(t, new(gooseSuite))
}

// GooseDbVersion is based on the schema of `goose_db_version` table.
type GooseDbVersion struct {
	ID        int64 `gorm:"column:id"`
	VersionID int64 `gorm:"column:version_id"`
	Applied   bool  `gorm:"column:is_applied"`
}

func (s *gooseSuite) TestUp() {
	migration := NewGooseMigration(DriverMysql, s.migrationDir(), s.mysqlURL())
	s.Require().NoError(migration.Up())
	s.Require().NoError(migration.Close())

	versions := []GooseDbVersion{}
	s.Require().NoError(s.db.Order("id ASC").Table("goose_db_version").Limit(2).Find(&versions).Error)
	s.Require().Equal(2, len(versions))
}
