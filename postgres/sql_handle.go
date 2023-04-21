package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sterrasi/pinion/app"
	"github.com/sterrasi/pinion/db"
)

// db.SqlHandle interface implementation
type sqlHandleImpl struct {
	tx pgx.Tx
}

// Exec executes the given SQL statement
func (sh *sqlHandleImpl) Exec(ctx context.Context, sql string, args ...any) (*db.ExecResult, app.Error) {
	tag, err := sh.tx.Exec(ctx, sql, args)
	if err != nil {
		return nil, handlePgxError(err, &statementDescriptor{
			operation: "execute",
			sql:       sql})
	}
	return &db.ExecResult{
		RowsAffected: uint(tag.RowsAffected()),
		Id:           tag.String(),
	}, nil
}

// Query executes the given sql query returning the matching rows
func (sh *sqlHandleImpl) Query(ctx context.Context, sql string, args ...any) (db.RowIterator, app.Error) {
	rows, err := sh.tx.Query(ctx, sql, args)
	if err != nil {
		return nil, handlePgxError(err, &statementDescriptor{
			operation: "query",
			sql:       sql})
	}
	return &pgxRowsWrapper{
		rows: rows,
		desc: &statementDescriptor{
			operation: "query",
			sql:       sql,
		},
	}, nil
}

// QueryRow executes a sql statement expecting only one row to be selected
func (sh *sqlHandleImpl) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	row := sh.tx.QueryRow(ctx, sql, args)
	return &pgxRowWrapper{
		row: row,
		desc: &statementDescriptor{
			operation: "single-row query",
			sql:       sql,
		},
	}
}

// pgxRowWrapper implements a db.Row
type pgxRowWrapper struct {
	row  pgx.Row
	desc *statementDescriptor
}

func (rs *pgxRowWrapper) Scan(dest ...any) app.Error {
	err := rs.row.Scan(dest)
	if err != nil {
		return handlePgxError(err, rs.desc)
	}
	return nil
}

// pgxRowsWrapper implements a db.RowIterator
type pgxRowsWrapper struct {
	rows pgx.Rows
	desc *statementDescriptor
}

func (rs *pgxRowsWrapper) Close() {
	rs.rows.Close()
}

func (rs *pgxRowsWrapper) Err() app.Error {
	err := rs.rows.Err()
	if err != nil {
		return handlePgxError(err, rs.desc)
	}
	return nil
}

func (rs *pgxRowsWrapper) Next() bool {
	return rs.rows.Next()
}

func (rs *pgxRowsWrapper) Scan(dest ...any) app.Error {
	err := rs.rows.Scan(dest)
	if err != nil {
		return handlePgxError(err, rs.desc)
	}
	return nil
}

func (rs *pgxRowsWrapper) Values() ([]any, app.Error) {
	values, err := rs.rows.Values()
	if err != nil {
		return nil, handlePgxError(err, rs.desc)
	}
	return values, nil
}
