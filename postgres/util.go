package postgres

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sterrasi/pinion/app"
	"github.com/sterrasi/pinion/db"
	"strings"
)

// postgres Error Codes
const duplicateKeyPgError = "23505"     // duplicate key constraint violation error (record already exists)
const connectionExceptionPgClass = "08" // Class 08 - Connection Exception

// isDuplicateKeyError checks the given error to see if it is a Postgres based error related to a duplicate
// key violation
func isDuplicateKeyError(err error) bool {
	switch err.(type) {
	case *pgconn.PgError:
		if err.(*pgconn.PgError).Code == duplicateKeyPgError {
			return true
		}
		return false

	default:
		return false
	}
}

// toPostgresError casts the given error to a postgres based error, otherwise returns an IllegalArgumentError
func toPostgresError(err error) (*pgconn.PgError, error) {
	switch err.(type) {
	case *pgconn.PgError:
		return err.(*pgconn.PgError), nil
	default:
		return nil, app.NewIllegalArgumentError("Expecting a Postgres error but got %T", err)
	}
}

// IsConnectionClassError returns true if the given error is Postgres based and if it is connection related
func IsConnectionClassError(err error) bool {
	switch err.(type) {
	case *pgconn.PgError:
		if strings.HasPrefix(err.(*pgconn.PgError).Code, connectionExceptionPgClass) {
			return true
		}
		return false

	default:
		return false
	}
}

type statementDescriptor struct {
	sql       string
	operation string
	context   string
}

func (sd *statementDescriptor) decorate(builder *app.ErrorBuilder) *app.ErrorBuilder {

	if sd.sql != "" {
		builder.Str("sql", sd.sql)
	}
	if sd.operation != "" {
		builder.Str("operation", sd.operation)
	}
	if sd.context != "" {
		builder.Context(sd.context)
	}
	return builder
}

func handlePgxError(err error, desc *statementDescriptor) app.Error {
	return handlePgxErrorWithCause(err, nil, desc)
}

func handlePgxErrorWithCause(err error, cause error, desc *statementDescriptor) app.Error {
	if err == nil {
		return nil
	}

	// check to see if the record already exists
	if isDuplicateKeyError(err) {
		errBuilder := app.BuildAlreadyExistsError().Cause(cause)
		pgError, e := toPostgresError(cause)
		if e == nil {
			errBuilder.
				Str("constraintName", pgError.ConstraintName).
				Str("postgresCode", pgError.Code).
				Str("postgresColumn", pgError.ColumnName)
		}

		return desc.decorate(errBuilder).Msgf("%s failed due to duplicate key", desc.operation)
	}

	// check for a database connection error
	if IsConnectionClassError(cause) {
		errBuilder := app.BuildSvcUnavailableError().Cause(err)
		pgError, e := toPostgresError(cause)
		if e == nil {
			errBuilder.Str("postgresCode", pgError.Code)
		}

		return desc.decorate(errBuilder).Msg("Database connection error")
	}

	// default to SQL exception
	errBuilder := db.BuildSqlError().Cause(err)
	pgError, e := toPostgresError(cause)
	if e == nil {
		errBuilder.Str("postgresCode", pgError.Code)
	}

	return desc.decorate(errBuilder).Msgf("SQL error during '%s' operation", desc.operation)
}
