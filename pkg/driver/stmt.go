package driver

import (
	"context"
	"database/sql/driver"
)

type SQLiteOGStmt struct {
	c        *SQLiteOGConn
	sql      string
	numInput int
}

func (s *SQLiteOGStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.c.ExecContext(ctx, s.sql, args)
}

func (s *SQLiteOGStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return s.c.QueryContext(ctx, s.sql, args)
}

func (s *SQLiteOGStmt) Close() error {
	return nil
}

func (s *SQLiteOGStmt) NumInput() int {
	return s.numInput
}

func (s *SQLiteOGStmt) Exec(args []driver.Value) (driver.Result, error) {
	namedValues := make([]driver.NamedValue, len(args))
	for k, v := range args {
		namedValues[0] = driver.NamedValue{
			Name:    "",
			Ordinal: k,
			Value:   v,
		}
	}
	return s.ExecContext(context.Background(), namedValues)
}

func (s *SQLiteOGStmt) Query(args []driver.Value) (driver.Rows, error) {
	namedValues := make([]driver.NamedValue, len(args))
	for k, v := range args {
		namedValues[0] = driver.NamedValue{
			Name:    "",
			Ordinal: k,
			Value:   v,
		}
	}
	return s.QueryContext(context.Background(), namedValues)
}
