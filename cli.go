package main

import (
	"fmt"

	"github.com/jtarchie/sqlettus/db"
	"github.com/jtarchie/sqlettus/handler"
	"github.com/jtarchie/sqlettus/tcp"
)

type CLI struct {
	Port     uint64 `default:"6379"    help:"port to listen on"`
	Filename string `default:"test.db" help:"filename to store database"`
}

func (c *CLI) Run() error {
	server := tcp.NewServer(c.Port)

	client, err := db.NewClient("test.db")
	if err != nil {
		return fmt.Errorf("could not start db client: %w", err)
	}
	defer client.Close()

	err = server.Listen(handler.New(client))
	if err != nil {
		return fmt.Errorf("could not listen for server: %w", err)
	}

	return nil
}
