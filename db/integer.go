package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jtarchie/sqlettus/db/writers"
)

func (c *Client) AddInt(ctx context.Context, name string, value int) (int, error) {
	intValue, err := c.writers.AddInt(ctx, writers.AddIntParams{
		Name:  name,
		Value: strconv.Itoa(value),
	})
	if err != nil {
		return 0, fmt.Errorf("could not ADDINT: %w", err)
	}

	return int(intValue), nil
}
