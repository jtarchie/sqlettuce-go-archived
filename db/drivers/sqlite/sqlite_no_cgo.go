//go:build !cgo
// +build !cgo

package sqlite

import (
	_ "modernc.org/sqlite"
)

const driverName = "sqlite"
