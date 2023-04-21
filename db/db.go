package db

import (
	"context"
	"github.com/sterrasi/pinion/app"
)

type Scanner interface {
	// Scan reads the values from the current row into dest values positionally.
	// dest can include pointers to core types, values implementing the Scanner
	// interface, and nil. nil will skip the value entirely. It is an error to
	// call Scan without first calling Next() and checking that it returned true.
	Scan(dest ...any) app.Error
}

type InsertStatement struct {
	SQL  string
	Name string
}

type MapperFn[M any] func(scanner Scanner, model *M) app.Error

type QueryStatement[M any] struct {
	SQL    string
	Name   string
	Mapper MapperFn[M]
}

type ExecResult struct {
	Id           string
	RowsAffected uint
}

type RowIterator interface {
	Scanner

	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() app.Error

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	// Values returns the decoded row values. As with Scan(), it is an error to
	// call Values without first calling Next() and checking that it returned
	// true.
	Values() ([]any, app.Error)
}

type Row interface {
	Scanner
}

type SqlHandle interface {
	Exec(ctx context.Context, sql string, args ...any) (*ExecResult, app.Error)
	Query(ctx context.Context, sql string, args ...any) (RowIterator, app.Error)
	QueryRow(ctx context.Context, sql string, args ...any) Row // no error?
}

type DatabaseHandle interface {
	SqlHandle
	ExecFile(filePath string) app.Error
	Insert(ctx context.Context, stmt *InsertStatement, args ...any) app.Error
}

type TransactionFn func(handle DatabaseHandle) app.Error
type DB interface {
	Close()
	ReadTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	WriteTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	WriteSerializableTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	Transaction(ctx context.Context, txOptions *TransactionOptions, tnFn TransactionFn) app.Error
}

// TxWriteOptions Transaction options for DB.WriteTransaction
var TxWriteOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: NotDeferrable,  // check constraints upfront
	IsoLevel:       RepeatableRead, // no phantom reads allowed in postgres
}

// TxSerializableWriteOptions Transaction options for DB.WriteSerializableTransaction
var TxSerializableWriteOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: NotDeferrable, // check constraints upfront
	IsoLevel:       Serializable,
}

// TxReadOptions Transaction options for DB.WriteSerializableTransaction
var TxReadOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: Deferrable,     // deffer constraint checks
	IsoLevel:       RepeatableRead, // removes non-repeatable reads
}

// ExistsValue can be used in a SingleRowQueryStatement to extract if a model exists
type ExistsValue struct {
	Exists bool
}
