package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jtarchie/sqlettus/db/readers"
	"github.com/jtarchie/sqlettus/db/writers"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Client struct {
	db *sql.DB

	readers *readers.Queries
	writers *writers.Queries
}

func NewClient(filename string) (*Client, error) {
	dsn := filename
	if !strings.Contains(dsn, "memory") {
		dsn = fmt.Sprintf("%s?_mutex=no&cache=share&_txlock=immediate&_locking=exclusive", filename)
	}

	writerDB, err := sql.Open(dbDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open writer db: %w", err)
	}

	writerDB.SetMaxOpenConns(1)

	// https://www.sqlite.org/compile.html
	// https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	_, err = writerDB.Exec(`
		PRAGMA busy_timeout = 5000;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA wal_autocheckpoint = 0;
		PRAGMA temp_store = memory;
		-- PRAGMA mmap_size = 30000000000;
		PRAGMA incremental_vacuum;
	`)
	if err != nil {
		return nil, fmt.Errorf("could not setup PRAGMA: %w", err)
	}

	migrationsFS, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not wrap migrations: %w", err)
	}

	driver, err := sqlite3.WithInstance(writerDB, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not wrap driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance(
		"iofs", migrationsFS,
		"ql", driver,
	)
	if err != nil {
		return nil, fmt.Errorf("could not setup migrations: %w", err)
	}

	err = migrator.Up()
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("could not migrate up: %w", err)
	}

	return &Client{
		db:      writerDB,
		readers: readers.New(writerDB),
		writers: writers.New(writerDB),
	}, nil
}

func (c *Client) Close() {
	c.readers.Close()
	c.writers.Close()

	c.db.Close()
}
