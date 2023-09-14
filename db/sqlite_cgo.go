//go:build cgo
// +build cgo

package db

import (
	// load the cgo version of sqlite.
	_ "github.com/mattn/go-sqlite3"
)

const dbDriverName = "sqlite3"
