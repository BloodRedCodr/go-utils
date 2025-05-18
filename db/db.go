package db

import (
	"crypto/tls"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-utils/auth"
	"go-utils/logger"
)

func Connect(
	logger *logger.Logger,
	p12CertPath string, p12Pwd string,
	host string, port string,
	username string, password string,
	dbName string,
) *gorm.DB {
	cert, caCertPool := auth.GetCertsFromP12(logger, p12CertPath, p12Pwd)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   host,
		MinVersion:   tls.VersionTLS12,
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full",
		host, port, username, password, dbName)

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		logger.Fatal("Failed to parse pgx config: %v", err)
	}
	pgxConfig.TLSConfig = tlsConfig

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: stdlib.OpenDB(*pgxConfig),
	}), &gorm.Config{})

	if err != nil {
		logger.Fatal("Failed to connect to the database: %v", err)
	}
	logger.Info("Database connection established successfully")

	return db
}

func RunMigrations(
	gormDB *gorm.DB, logger *logger.Logger,
	dbName string, migrationPath string,
) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("Failed to get raw DB: %v", err)
	}

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		logger.Fatal("Failed to get database driver for migration: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, dbName, driver)
	if err != nil {
		logger.Fatal("Failed to connect database for migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Failed to migrate the database: %v", err)
	}

	logger.Info("Database migration completed successfully")
}
