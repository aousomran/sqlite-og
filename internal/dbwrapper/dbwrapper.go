package dbwrapper

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aousomran/sqlite-og/internal/callback"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	pb "github.com/aousomran/sqlite-og/gen/proto"
)

const DefaultDBName = "test"
const defaultDriverName = "sqlite-og"

type Columns []string
type RowFields *[]string
type Rows []RowFields

type callbackFunction func(args ...interface{}) string

func makeCallbackFunc(functionName string, channels *callback.CallbackChannels) callbackFunction {
	return func(args ...interface{}) string {
		slog.Debug("got invocation from DB", "func_name", functionName, "args", args)
		args2 := make([]string, 0)
		for _, v := range args {
			args2 = append(args2, fmt.Sprintf("%v", v))
		}
		channels.ChanSend <- &pb.Invoke{
			FunctionName: functionName,
			Args:         args2,
		}
		result := <-channels.ChanReceive
		slog.Debug("received result sending back to DB", "result", result.GetResult())
		if len(result.GetResult()) < 1 {
			return ""
		}
		return result.GetResult()[0]
	}
}

func registerDriver(driverName string, functions []string, channels *callback.CallbackChannels) {
	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			for _, name := range functions {
				slog.Debug("registering functions", "names", functions)
				err := conn.RegisterFunc(name, makeCallbackFunc(name, channels), true)
				if err != nil {
					slog.Error("unable to register function", "name", name, "error", err.Error())
					return err
				}
			}
			return nil
		},
	})
	return
}

func normalizeDBName(name string) string {
	// allow in memory database
	if name == ":memory:" {
		return name
	}

	if name == "" {
		slog.Warn("got empty database name, using default", "dbname", DefaultDBName)
		return DefaultDBName
	}
	name = strings.TrimSpace(strings.TrimSuffix(name, ".db"))
	return fmt.Sprintf("%s.db", name)

}

type DBWrapper struct {
	Name     string
	Database *sql.DB
	Channels *callback.CallbackChannels
}

func New(dbname, id string, functions []string, channels *callback.CallbackChannels) *DBWrapper {
	// TODO: pass context to this function
	dbname = normalizeDBName(dbname)
	registerDriver(id, functions, channels)
	slog.Info("registered drivers", "names", sql.Drivers())
	return &DBWrapper{
		Name:     dbname,
		Channels: channels,
	}
}

func (w *DBWrapper) Open(driverName string) error {
	if w.Database == nil {
		db, err := sql.Open(driverName, w.Name)
		if err != nil {
			return err
		}

		w.Database = db
	}
	return nil
}

func (w *DBWrapper) Close() error {
	if w.Database != nil {
		err := w.Database.Close()
		if err != nil {
			return err
		}
	}
	w.Database = nil
	return nil
}

func (w *DBWrapper) Query(ctx context.Context, sql string, params ...interface{}) ([]string, []string, []*pb.Row, error) {
	if w.Database == nil {
		return nil, nil, nil, fmt.Errorf("connection is closed")
	}

	rows, err := w.Database.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, nil, nil, err
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			slog.ErrorContext(ctx, "unable to close rows", "error", errClose)
		}
	}()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, nil, err
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, nil, err
	}

	colTypes2 := make([]string, len(colTypes))
	for k, v := range colTypes {
		colTypes2[k] = v.DatabaseTypeName()
	}

	pbRows := make([]*pb.Row, 0)
	for rows.Next() {
		r, err := rowToStringSlice(cols, rows)
		if err != nil {
			slog.ErrorCtx(ctx, "unable to fetch next row", "error", err)
			return nil, nil, nil, err
		}
		pbRows = append(pbRows, &pb.Row{Fields: r})
	}

	return cols, colTypes2, pbRows, nil
}

func (w *DBWrapper) Execute(ctx context.Context, sql string, params ...interface{}) (insertId int64, affected int64, err error) {
	if w.Database == nil {
		err = fmt.Errorf("connection is closed")
		return
	}

	result, err := w.Database.ExecContext(ctx, sql, params...)
	if err != nil {
		return
	}

	insertId, err = result.LastInsertId()
	if err != nil {
		return
	}

	affected, err = result.RowsAffected()
	return
}

func rowToStringSlice(columnNames []string, rows *sql.Rows) ([]string, error) {
	lenCN := len(columnNames)
	ret := make([]string, lenCN)

	columnPointers := make([]interface{}, lenCN)
	for i := 0; i < lenCN; i++ {
		columnPointers[i] = new(sql.RawBytes)
	}

	if err := rows.Scan(columnPointers...); err != nil {
		return nil, err
	}

	for i := 0; i < lenCN; i++ {
		if rb, ok := columnPointers[i].(*sql.RawBytes); ok {
			ret[i] = string(*rb)
		} else {
			return nil, fmt.Errorf("cannot convert index %d column %s to type *sql.RawBytes", i, columnNames[i])
		}
	}

	return ret, nil
}
