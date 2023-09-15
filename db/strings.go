package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jtarchie/sqlettus/db/readers"
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

func (c *Client) Delete(ctx context.Context, name string) (*string, error) {
	value, err := c.writers.Delete(ctx, name)

	//nolint:nilnil
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("could not DELETE: %w", err)
	}

	return &value, nil
}

func (c *Client) Append(ctx context.Context, name, value string) (int, error) {
	length, err := c.writers.Append(ctx, writers.AppendParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return 0, fmt.Errorf("could not APPEND: %w", err)
	}

	return int(length.Int64), nil
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
