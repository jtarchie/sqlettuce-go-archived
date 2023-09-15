//go:build cgo
// +build cgo

package sqlite

import (
	// load the cgo version of sqlite.
	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"
