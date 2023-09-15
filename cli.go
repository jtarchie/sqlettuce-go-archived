package main

import (
	"context"
	"fmt"

	"github.com/jtarchie/sqlettus/db"
	"github.com/jtarchie/sqlettus/handler"
	"github.com/jtarchie/sqlettus/tcp"
)

type CLI struct {
	Port     uint   `default:"6379"    help:"port to listen on"`
	Filename string `default:"test.db" help:"filename to store database"`
	Workers  uint   `default:"100"     help:"number of workers to run"`
}

func (c *CLI) Run() error {
	ctx := context.TODO()

	client, err := db.NewClient(c.Filename)
	if err != nil {
		return fmt.Errorf("could not start db client: %w", err)
	}
	defer client.Close()

	server, err := tcp.NewServer(ctx, c.Port, c.Workers)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	err = server.Listen(ctx, handler.New(client))
	if err != nil {
		return fmt.Errorf("could not listen for server: %w", err)
	}

	return nil
}
