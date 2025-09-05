// Package db represents API for database operations.
package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zeta-chain/node/zetaclient/types"
)

// SqliteInMemory is a special string to use in-memory database.
// @see https://www.sqlite.org/inmemorydb.html
const SqliteInMemory = ":memory:"

// read/write/execute for user
// read/write for group
const dirCreationMode = 0o750

var (
	defaultGormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	migrationEntities = []any{
		&types.LastBlockSQLType{},
		&types.TransactionSQLType{},
		&types.ReceiptSQLType{},
		&types.TransactionResultSQLType{},
		&types.OutboundHashSQLType{},
		&types.LastTransactionSQLType{},
		&types.AuxStringSQLType{},
	}
)

// DB database.
type DB struct {
	db *gorm.DB
}

// NewFromSqlite creates a new instance of DB based on SQLite database.
func NewFromSqlite(directory, dbName string, migrate bool) (*DB, error) {
	path, err := ensurePath(directory, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to ensure database path")
	}

	return New(sqlite.Open(path), migrate)
}

// NewFromSqliteInMemory creates a new instance of DB based on SQLite in-memory database.
func NewFromSqliteInMemory(migrate bool) (*DB, error) {
	return NewFromSqlite(SqliteInMemory, "", migrate)
}

// New creates a new instance of DB.
func New(dial gorm.Dialector, migrate bool) (*DB, error) {
	// open db
	db, err := gorm.Open(dial, defaultGormConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open gorm database")
	}

	if migrate {
		if err := db.AutoMigrate(migrationEntities...); err != nil {
			return nil, errors.Wrap(err, "unable to migrate database")
		}
	}

	return &DB{db}, nil
}

// Client returns the underlying gorm database.
func (db *DB) Client() *gorm.DB {
	return db.db
}

// Close closes the database.
func (db *DB) Close() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return errors.Wrap(err, "unable to get underlying sql.DB")
	}

	if err := sqlDB.Close(); err != nil {
		return errors.Wrap(err, "unable to close sql.DB")
	}

	return nil
}

func ensurePath(directory, dbName string) (string, error) {
	// pass in-memory database as is
	if strings.Contains(directory, SqliteInMemory) {
		return directory, nil
	}

	_, err := os.Stat(directory)
	switch {
	case os.IsNotExist(err):
		if err := os.MkdirAll(directory, dirCreationMode); err != nil {
			return "", errors.Wrapf(err, "unable to create database path %q", directory)
		}
	case err != nil:
		return "", errors.Wrap(err, "unable to check database path")
	}

	return fmt.Sprintf("%s/%s", directory, dbName), nil
}
