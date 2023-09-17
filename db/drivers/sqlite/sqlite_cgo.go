//go:build cgo
// +build cgo

package sqlite

import (
	// load the cgo version of sqlite.
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3_custom"

func init() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("json_array_insert", jsonArrayInsert, true); err != nil {
				return fmt.Errorf("could not register json_array_insert: %w", err)
			}

			return nil
		},
	})
}

var ErrNotArray = errors.New("not an array")

func jsonArrayInsert(array string, pivot, value string, offset int) (string, error) {
	if array[0] != '[' {
		return "", ErrNotArray
	}

	var elements []string

	err := json.Unmarshal([]byte(array), &elements)
	if err != nil {
		return "", fmt.Errorf("could not parse JSON: %w", err)
	}

	for index, element := range elements {
		if pivot == element {
			if offset < 0 {
				elements = append(elements[:index], append([]string{value}, elements[index:]...)...)
			} else {
				elements = append(elements[:index+1], append([]string{value}, elements[index+1:]...)...)
			}

			break
		}
	}

	contents, err := json.Marshal(elements)
	if err != nil {
		return "", fmt.Errorf("could not marshal JSON: %w", err)
	}

	return string(contents), nil
}
