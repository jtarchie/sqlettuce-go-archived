package db

import (
	"context"
	"fmt"
)

func (c *Client) FlushAll(ctx context.Context) error {
	err := c.writers.FlushAll(ctx)
	if err != nil {
		return fmt.Errorf("could not flush all: %w", err)
	}

	return nil
}
