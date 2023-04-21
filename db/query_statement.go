package db

import (
	"context"
	"github.com/sterrasi/pinion/app"
	"github.com/sterrasi/pinion/logger"
)

func (q *QueryStatement[M]) QueryRow(ctx context.Context, handle SqlHandle, args ...any) (*M, app.Error) {

	logger.Debug().
		Str("queryName", q.Name).
		Msg("Executing single row query")

	row := handle.QueryRow(ctx, q.SQL, args...)
	model := new(M)
	err := q.Mapper(row, model)
	if err != nil {
		err.SetContext(q.Name)
		return nil, err
	}

	logger.Debug().
		Str("queryName", q.Name).
		Msg("Successfully executed query")

	return model, nil
}

func (q *QueryStatement[M]) Query(ctx context.Context, handle SqlHandle, args ...any) ([]*M, app.Error) {

	logger.Debug().
		Str("queryName", q.Name).
		Msg("Executing multi-row query")

	rows, appErr := handle.Query(ctx, q.SQL, args...)
	if appErr != nil {
		appErr.SetContext(q.Name)
		return nil, appErr
	}
	defer rows.Close()
	results := make([]*M, 0, 1)

	for rows.Next() {

		model := new(M)
		if appErr = q.Mapper(rows, model); appErr != nil {
			appErr.SetContext(q.Name)
			return nil, appErr
		}
		results = append(results, model)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if appErr = rows.Err(); appErr != nil {
		appErr.SetContext(q.Name)
		return nil, appErr
	}
	return results, nil
}
