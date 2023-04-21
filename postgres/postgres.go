package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sterrasi/pinion/app"
	"github.com/sterrasi/pinion/db"
	"github.com/sterrasi/pinion/logger"
)

// pgDb is a postgres specific (pgx) DB interface
type pgDb struct {
	dbHandleImpl
	pool   *pgxpool.Pool
	Config *db.DbConfig
	url    string
}

// NewPostgresDb connects to a postgres database described in the given db.DbConfig
func NewPostgresDb(cfg *db.DbConfig) (db.DB, app.Error) {

	pg := &pgDb{Config: cfg}
	pg.url = fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.DbName)

	dbPool, err := pgxpool.New(context.Background(), pg.url)
	if err != nil {
		return nil, app.BuildSysConfigError().
			Cause(err).
			Str("url", pg.url).
			Msg("Error creating database pool")
	}
	pg.pool = dbPool

	// Ping the database to make sure the connection is valid
	if err = pg.pool.Ping(context.Background()); err != nil {
		return nil, app.BuildSvcUnavailableError().
			Cause(err).
			Str("url", pg.url).
			Msg("Error connecting to postgres")
	}

	var hiddenPassword string
	if len(cfg.Password) > 0 {
		hiddenPassword = "<yes>"
	} else {
		hiddenPassword = "<no>"
	}

	// log that the database was successfully connected to
	logger.Info().
		Str("user", cfg.User).
		Str("password", hiddenPassword).
		Str("host", cfg.Host).
		Str("schema", cfg.Schema).
		Str("database name", cfg.DbName).
		Str("url", pg.url).
		Msg("Connected to postgres")

	return pg, nil
}

// Close closes the connection
func (pg *pgDb) Close() {
	pg.pool.Close()
	logger.Info().
		Str("url", pg.url).
		Msg("Successfully shut down connection to postgres")
}

// ReadTransaction executes a transaction with db.TransactionOptions defaults for read-only
func (pg *pgDb) ReadTransaction(ctx context.Context, tnFn db.TransactionFn) app.Error {
	return pg.Transaction(ctx, db.TxReadOptions, tnFn)
}

// WriteTransaction executes a transaction with db.TransactionOptions defaults for read-write
func (pg *pgDb) WriteTransaction(ctx context.Context, tnFn db.TransactionFn) app.Error {
	return pg.Transaction(ctx, db.TxWriteOptions, tnFn)
}

// WriteSerializableTransaction executes a transaction with serializable isolation level
func (pg *pgDb) WriteSerializableTransaction(ctx context.Context, tnFn db.TransactionFn) app.Error {
	return pg.Transaction(ctx, db.TxSerializableWriteOptions, tnFn)
}

// Transaction executes a transaction with pgx.TxOptions defaults for read-only
func (pg *pgDb) Transaction(ctx context.Context, txOptions *db.TransactionOptions, tnFn db.TransactionFn) app.Error {

	pgxOpts, appErr := asPgxOptions(txOptions)
	if appErr != nil {
		return appErr
	}

	tx, err := pg.pool.BeginTx(ctx, *pgxOpts)
	if err != nil {
		return db.BuildDatabaseError().
			Cause(err).
			Msg("Error calling begin transaction")
	}

	appErr = tnFn(NewDatabaseHandle(tx))
	if appErr != nil {
		_ = tx.Rollback(ctx)
		return appErr
	}
	err = tx.Commit(ctx)
	if err != nil {
		return db.BuildSqlError().
			Cause(err).
			Msg("Error Committing transaction")
	}
	return nil
}

func asPgxOptions(stdOpts *db.TransactionOptions) (*pgx.TxOptions, app.Error) {
	pgxOpts := &pgx.TxOptions{}

	// convert Iso Level
	switch {
	case stdOpts.IsoLevel == db.Serializable:
		pgxOpts.IsoLevel = pgx.Serializable
	case stdOpts.IsoLevel == db.RepeatableRead:
		pgxOpts.IsoLevel = pgx.RepeatableRead
	case stdOpts.IsoLevel == db.ReadCommitted:
		pgxOpts.IsoLevel = pgx.ReadCommitted
	case stdOpts.IsoLevel == db.ReadUncommitted:
		pgxOpts.IsoLevel = pgx.ReadUncommitted
	default:
		return nil, app.BuildIllegalArgumentError().
			Str("value", fmt.Sprintf("%v", stdOpts.IsoLevel)).
			Msg("Invalid Transaction Options Iso Level")
	}

	// convert Deferrable Mode
	switch {
	case stdOpts.DeferrableMode == db.Deferrable:
		pgxOpts.DeferrableMode = pgx.Deferrable
	case stdOpts.DeferrableMode == db.NotDeferrable:
		pgxOpts.DeferrableMode = pgx.NotDeferrable
	default:
		return nil, app.BuildIllegalArgumentError().
			Str("value", fmt.Sprintf("%v", stdOpts.DeferrableMode)).
			Msg("Invalid Transaction Options Deferrable Mode")
	}

	// convert Access Mode
	switch {
	case stdOpts.AccessMode == db.ReadWrite:
		pgxOpts.AccessMode = pgx.ReadWrite
	case stdOpts.AccessMode == db.ReadOnly:
		pgxOpts.AccessMode = pgx.ReadOnly
	default:
		return nil, app.BuildIllegalArgumentError().
			Str("value", fmt.Sprintf("%v", stdOpts.AccessMode)).
			Msg("Invalid Transaction Options Access Mode")

	}

	return pgxOpts, nil
}
