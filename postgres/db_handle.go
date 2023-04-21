package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sterrasi/pinion/app"
	"github.com/sterrasi/pinion/db"
	"github.com/sterrasi/pinion/logger"
	"os"
	"strconv"
)

// implementation struct for a db.DatabaseHandle
type dbHandleImpl struct {
	db.SqlHandle
}

// NewDatabaseHandle creates a db.DatabaseHandle from the given transaction
func NewDatabaseHandle(tx pgx.Tx) db.DatabaseHandle {
	return &dbHandleImpl{&sqlHandleImpl{tx: tx}}
}

// ExecFile will execute the given SQL file
func (dh *dbHandleImpl) ExecFile(filePath string) app.Error {

	c, err := os.ReadFile(filePath)
	if err != nil {
		return app.BuildIOError().
			Str("file", filePath).
			Msg("Error reading file")
	}
	sql := string(c)
	_, appErr := dh.Exec(context.Background(), sql)
	if appErr != nil {
		return appErr
	}
	return nil
}

// Insert will execute an insert statement and expect the addition of a record as a result
func (dh *dbHandleImpl) Insert(ctx context.Context, stmt *db.InsertStatement, args ...any) app.Error {

	tag, err := dh.doInsert(ctx, stmt, args...)
	if err != nil {
		return err
	}
	// make sure that there was an affected row
	if tag.RowsAffected == 0 {
		return db.BuildSqlError().
			Context(stmt.Name).
			Msg("Record was not inserted")
	}
	return nil
}

// InsertOpt will execute an insert statement and return the number of rows affected
func (dh *dbHandleImpl) InsertOpt(ctx context.Context, stmt *db.InsertStatement, args ...any) (*uint64, app.Error) {

	tag, err := dh.doInsert(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}

	// make sure that there was an affected row
	affected := uint64(tag.RowsAffected)
	return &affected, nil
}

func (dh *dbHandleImpl) doInsert(ctx context.Context, stmt *db.InsertStatement, args ...any) (*db.ExecResult,
	app.Error) {

	logger.Debug().
		Str("statementName", stmt.Name).
		Msg("Executing Insert")

	tag, err := dh.Exec(ctx, stmt.SQL, args...)
	if err != nil {
		err.SetContext(stmt.Name)
		return nil, err
	}

	logger.Debug().
		Str("statementName", stmt.Name).
		Str("rowsAffected", strconv.FormatUint(uint64(tag.RowsAffected), 10)).
		Msg("Executed Insert")

	return tag, nil
}
