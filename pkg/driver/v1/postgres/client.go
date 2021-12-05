package postgres

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	DB     *sqlx.DB
	dBName string
}

func NewClient(uri, dbName string) (Client, error) {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return Client{}, err
	}
	return Client{DB: db, dBName: dbName}, nil
}

func (c Client) Migrate(migratePath string) error {
	driver, err := postgres.WithInstance(c.DB.DB, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    c.dBName,
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", migratePath),
		c.dBName, driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Errorf("error when migration up: %v", err)
	}
	return nil
}
