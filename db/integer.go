package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jtarchie/sqlettus/db/drivers/sqlite/writers"
)

func (c *Client) AddInt(ctx context.Context, name string, value int64) (int64, error) {
	intValue, err := c.writers.AddInt(ctx, writers.AddIntParams{
		Name:  name,
		Value: strconv.FormatInt(value, 10),
	})
	if err != nil {
		return 0, fmt.Errorf("could not ADDINT: %w", err)
	}

	return intValue, nil
}
