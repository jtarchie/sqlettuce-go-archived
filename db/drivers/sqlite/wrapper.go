//nolint:ireturn
package sqlite

import (
	"database/sql"

	"github.com/jtarchie/sqlettuce/db/drivers/sqlite/batch"
	"github.com/jtarchie/sqlettuce/db/drivers/sqlite/readers"
	"github.com/jtarchie/sqlettuce/db/drivers/sqlite/writers"
)

// Done to wrap the WithTX function
// https://github.com/sqlc-dev/sqlc/issues/383#issuecomment-613083090

type Reader interface {
	readers.Querier
	Close() error
	WithTx(tx *sql.Tx) Reader
}

type Readers struct {
	*readers.Queries
}

func (q *Readers) WithTx(tx *sql.Tx) Reader {
	return &Readers{
		Queries: q.Queries.WithTx(tx),
	}
}

type Writer interface {
	writers.Querier
	Close() error
	WithTx(tx *sql.Tx) Writer
}

type Writers struct {
	*writers.Queries
}

func (q *Writers) WithTx(tx *sql.Tx) Writer {
	return &Writers{
		Queries: q.Queries.WithTx(tx),
	}
}

type Batcher interface {
	batch.Querier
	WithTx(tx *sql.Tx) Batcher
}

type Batches struct {
	*batch.Queries
}

func (q *Batches) WithTx(tx *sql.Tx) Batcher {
	return &Batches{
		Queries: q.Queries.WithTx(tx),
	}
}

type Driver struct {
	DB *sql.DB

	Readers Reader
	Writers Writer
	Batcher Batcher
}
