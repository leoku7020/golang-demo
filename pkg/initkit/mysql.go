package initkit

import (
	"os"
	"strconv"
	"time"

	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"gorm.io/gorm"

	"demo/pkg/envkit"
	"demo/pkg/logger"
	sqlMig "demo/pkg/migration/sql"
	"demo/pkg/mysqlkit"
)

func NewGormDB() *gorm.DB {
	dsn := os.Getenv("DSN")
	maxOpenConns, _ := strconv.ParseInt(os.Getenv("MAX_OPEN_CONNS"), 10, 64)
	connMaxLifeMinutes, _ := strconv.ParseInt(os.Getenv("CONN_MAX_LIFE_MINUTES"), 10, 64)
	logger.Debug("GORM preparing", logger.WithFields(logger.Fields{
		"dsn":                dsn,
		"max-open-conns":     maxOpenConns,
		"conn-max-life-mins": connMaxLifeMinutes,
	}))

	cfg, err := mysqlkit.NewConfig(dsn)
	if err != nil {
		logger.Fatal("mysqlkit.NewConfig failed", logger.WithError(err))
	}

	db, err := mysqlkit.NewGORM(
		cfg.FormatDSN(),
		telemetry.NewNrTracer(
			cfg.DBName, // db name
			cfg.Addr,   // Host is the name of the server hosting the datastore.
			"MySQL",    // product name: fixed string defined by new relic
		),
	)
	if err != nil {
		logger.Fatal("mysqlkit.NewGORM failed", logger.WithError(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("db.DB", logger.WithError(err))
	}
	sqlDB.SetMaxOpenConns(int(maxOpenConns))
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifeMinutes) * time.Minute)

	logger.Info("GORM done", logger.WithFields(logger.Fields{
		"max-open-conns":     maxOpenConns,
		"conn-max-life-mins": connMaxLifeMinutes,
	}))

	return db
}

func ExecMySQLMigration() {
	// only allow to migrate schemas in development env or in test case.
	// otherwise, perform migration process in presync step in staging and production
	if envkit.Namespace() != envkit.EnvDevelopment {
		return
	}
	dsn := os.Getenv("DSN")
	migrationDir := os.Getenv("MIGRATION_DIR")

	if migrationDir == "" {
		logger.Fatal("MigrationDir should not be empty")
	}
	logger.Debug("Migration preparing", logger.WithFields(logger.Fields{
		"dsn": dsn,
		"dir": migrationDir,
	}))

	cfg, err := mysqlkit.NewConfig(dsn)
	if err != nil {
		logger.Fatal("mysqlkit.NewConfig failed", logger.WithError(err))
	}

	migration := sqlMig.NewGooseMigration(sqlMig.DriverMysql, migrationDir, cfg.FormatDSN())
	if err := migration.Up(); err != nil {
		logger.Fatal("migration.Up failed", logger.WithError(err))
	}
	migration.Close()

	logger.Info("Migration done", logger.WithFields(logger.Fields{
		"dir": migrationDir,
	}))
}
