package db

import (
	"context"
	"fmt"

	"github.com/jtarchie/sqlettus/db/writers"
)

func (c *Client) Set(name, value string) error {
	err := c.writers.Set(context.Background(), writers.SetParams{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("could not SET: %w", err)
	}

	return nil
}

func (c *Client) Get(name string) (*string, error) {
	value, err := c.readers.Get(context.Background(), name)
	if err != nil {
		return nil, fmt.Errorf("could not GET: %w", err)
	}

	return &value, nil
}
