//go:build !cgo
// +build !cgo

package db

import (
	_ "modernc.org/sqlite"
)

const dbDriverName = "sqlite"
