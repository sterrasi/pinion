package db

import (
	"context"
	"github.com/sterrasi/pinion/app"
)

type Scanner interface {
	Scan(dest ...any) error
}

type InsertStatement struct {
	SQL  string
	Name string
}

type MapperFn[M any] func(scanner Scanner, model *M) error

type QueryStatement[M any] struct {
	SQL    string
	Name   string
	Mapper MapperFn[M]
}

type StatusText struct {
	status string
}

type RowIterator interface {

	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() error

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	// Scan reads the values from the current row into dest values positionally.
	// dest can include pointers to core types, values implementing the Scanner
	// interface, and nil. nil will skip the value entirely. It is an error to
	// call Scan without first calling Next() and checking that it returned true.
	Scan(dest ...any) error

	// Values returns the decoded row values. As with Scan(), it is an error to
	// call Values without first calling Next() and checking that it returned
	// true.
	Values() ([]any, error)

	// RawValues returns the unparsed bytes of the row values. The returned data is only valid until the next Next
	// call or the Rows is closed.
	RawValues() [][]byte
}

type Row interface {
	// Scan works the same as Rows. with the following exceptions. If no
	// rows were found it returns ErrNoRows. If multiple rows are returned it
	// ignores all but the first.
	Scan(dest ...any) error
}

type SQLHandle interface {
	Exec(ctx context.Context, sql string, args ...any) (StatusText, error)
	Query(ctx context.Context, sql string, args ...any) (RowIterator, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row // no error?
}

type DbHandle interface {
	SQLHandle
	ExecFile(filePath string) app.Error
	Insert(ctx context.Context, stmt *InsertStatement, args ...any) app.Error
}

type TransactionFn func(handle DbHandle) app.Error
type DB interface {
	DbHandle
	Close()
	ReadTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	WriteTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	WriteSerializableTransaction(ctx context.Context, tnFn TransactionFn) app.Error
	Transaction(ctx context.Context, txOptions *TransactionOptions, tnFn TransactionFn) app.Error
}

// Transaction options for DB.WriteTransaction
var txWriteOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: NotDeferrable,  // check constraints upfront
	IsoLevel:       RepeatableRead, // no phantom reads allowed in postgres
}

// Transaction options for DB.WriteSerializableTransaction
var txSerializableWriteOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: NotDeferrable, // check constraints upfront
	IsoLevel:       Serializable,
}

// Transaction options for DB.WriteSerializableTransaction
var txReadOptions = &TransactionOptions{
	AccessMode:     ReadWrite,
	DeferrableMode: Deferrable,     // deffer constraint checks
	IsoLevel:       RepeatableRead, // removes non-repeatable reads
}

// ExistsValue can be used in a SingleRowQueryStatement to extract if a model exists
type ExistsValue struct {
	Exists bool
}

// postgres error codes
const duplicateKeyPgError = "23505"     // duplicate key constraint violation error (record already exists on insert/update)
const connectionExceptionPgClass = "08" // Class 08 - Connection Exception

//// IsDuplicateKeyError checks the given error to see if it is a Postgres based error related to a duplicate key
//// violation
//func IsDuplicateKeyError(err error) bool {
//	switch err.(type) {
//	case *pgconn.PgError:
//		if err.(*pgconn.PgError).Code == duplicateKeyPgError {
//			return true
//		}
//		return false
//
//	default:
//		return false
//	}
//}
