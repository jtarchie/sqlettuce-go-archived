package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettus/db/writers"
)

func (c *Client) Set(ctx context.Context, name, value string) error {
	err := c.writers.Set(ctx, writers.SetParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("could not SET: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, name string) (*string, error) {
	value, err := c.readers.Get(ctx, name)

	//nolint:nilnil
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not GET: %w", err)
	}

	return &value, nil
}

func (c *Client) Delete(ctx context.Context, name string) error {
	err := c.writers.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("could not DELETE: %w", err)
	}

	return nil
}
