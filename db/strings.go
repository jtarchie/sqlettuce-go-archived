package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/jtarchie/sqlettus/db/drivers/sqlite/readers"
	"github.com/jtarchie/sqlettus/db/drivers/sqlite/writers"
)

func (c *Client) Set(ctx context.Context, name, value string) error {
	err := c.writers.Set(ctx, &writers.SetParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("could not SET: %w", err)
	}

	return nil
}

func (c *Client) MSet(ctx context.Context, args ...string) error {
	transaction, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start MSET: %w", err)
	}
	//nolint:errcheck
	defer transaction.Rollback()

	queries := c.writers.WithTx(transaction)
	params := &writers.SetParams{}

	for index := 0; index < len(args); index += 2 {
		params.Name = args[index]
		params.Value = args[index+1]

		err := queries.Set(ctx, params)
		if err != nil {
			return fmt.Errorf("could not set MSET: %w", err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("could not MSET: %w", err)
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

func (c *Client) MGet(ctx context.Context, names ...string) ([]string, error) {
	results, err := c.batcher.Get(ctx, names)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not MGET: %w", err)
	}

	values := make([]string, len(names))

	for _, result := range results {
		index := slices.Index(names, result.Name)

		if index >= 0 {
			values[index] = result.Value
		}
	}

	return values, nil
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
	length, err := c.writers.AppendValue(ctx, &writers.AppendValueParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return 0, fmt.Errorf("could not APPEND: %w", err)
	}

	return length.Int64, nil
}

func (c *Client) Substr(ctx context.Context, name string, start, end int64) (string, error) {
	value, err := c.readers.Substr(ctx, &readers.SubstrParams{
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
