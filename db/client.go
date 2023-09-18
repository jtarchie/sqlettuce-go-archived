package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/jtarchie/sqlettuce/db/drivers/sqlite"
)

type Client struct {
	db *sql.DB

	readers sqlite.Reader
	writers sqlite.Writer
	batcher sqlite.Batcher
}

var ErrDriverNotFound = errors.New("could not find driver")

func NewClient(dsn string) (*Client, error) {
	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse DSN (%q): %w", dsn, err)
	}

	switch uri.Scheme {
	case "sqlite":
		driver, err := sqlite.New(uri.Host)
		if err != nil {
			return nil, fmt.Errorf("could not setup sqlite: %w", err)
		}

		return &Client{
			db:      driver.DB,
			readers: driver.Readers,
			writers: driver.Writers,
			batcher: driver.Batcher,
		}, nil
	default:
		return nil, fmt.Errorf("could not find a driver for %q: %w", uri.Scheme, ErrDriverNotFound)
	}
}

func (c *Client) Close() error {
	err := c.readers.Close()
	if err != nil {
		return fmt.Errorf("could not close readers: %w", err)
	}

	err = c.writers.Close()
	if err != nil {
		return fmt.Errorf("could not close writers: %w", err)
	}

	err = c.db.Close()
	if err != nil {
		return fmt.Errorf("could not close db: %w", err)
	}

	return nil
}
