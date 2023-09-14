package db

import (
	"context"
	"fmt"
)

func (c *Client) FlushAll() error {
	err := c.writers.FlushAll(context.Background())
	if err != nil {
		return fmt.Errorf("could not flush all: %w", err)
	}

	return nil
}
