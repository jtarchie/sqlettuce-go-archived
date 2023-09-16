package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jtarchie/sqlettus/db/drivers/sqlite/writers"
)

func (c *Client) AddFloat(ctx context.Context, name string, value float64) (float64, error) {
	newValue, err := c.writers.AddFloat(ctx, writers.AddFloatParams{
		Name:  name,
		Value: strconv.FormatFloat(value, 'f', 17, 64),
	})
	if err != nil {
		return 0, fmt.Errorf("could not ADDFLOAT: %w", err)
	}

	return newValue, nil
}
