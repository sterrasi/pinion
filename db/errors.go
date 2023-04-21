package db

import "github.com/sterrasi/pinion/app"

// DatabaseOperationErrorCode signifies a database level operation that failed
const DatabaseOperationErrorCode app.ErrorCode = 10

func BuildDatabaseError() *app.ErrorBuilder {
	return app.NewErrorBuilder(DatabaseOperationErrorCode, "database")
}
func NewDatabaseError(format string, args ...any) app.Error {
	return app.BuildValidationError().Msgf(format, args...)
}

// SQLErrorCode signifies a failure for the database to handle an SQL statement
const SQLErrorCode app.ErrorCode = 10

func BuildSqlError() *app.ErrorBuilder {
	return app.NewErrorBuilder(SQLErrorCode, "sql")
}
func NewSQLError(format string, args ...any) app.Error {
	return app.BuildValidationError().Msgf(format, args...)
}
