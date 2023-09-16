package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettus/db/drivers/sqlite/readers"
	"github.com/jtarchie/sqlettus/db/drivers/sqlite/writers"
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

func (c *Client) Get(ctx context.Context, name string) (string, bool, error) {
	value, err := c.readers.Get(ctx, name)

	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, fmt.Errorf("could not GET: %w", err)
	}

	return value, true, nil
}

func (c *Client) Delete(ctx context.Context, names ...string) ([]string, bool, error) {
	values, err := c.batcher.Delete(ctx, names)

	if errors.Is(err, sql.ErrNoRows) || len(values) == 0 {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("could not DELETE: %w", err)
	}

	return values, true, nil
}

func (c *Client) Append(ctx context.Context, name, value string) (int64, error) {
	length, err := c.writers.Append(ctx, writers.AppendParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return 0, fmt.Errorf("could not APPEND: %w", err)
	}

	return length.Int64, nil
}

func (c *Client) Substr(ctx context.Context, name string, start, end int64) (string, error) {
	value, err := c.readers.Substr(ctx, readers.SubstrParams{
		Name:  name,
		Start: start,
		End:   end,
	})

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("could not SUBSTR: %w", err)
	}

	return value, nil
}
