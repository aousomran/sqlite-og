package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	pb "github.com/aousomran/sqlite-og/gen/proto"
	"google.golang.org/grpc"
)

type SQLiteOGConn struct {
	ID       string
	GRPCConn *grpc.ClientConn
	OGClient pb.SqliteOGClient
}

func (c *SQLiteOGConn) Prepare(query string) (driver.Stmt, error) {
	if query == "" {
		return nil, errors.New("query is empty")
	}
	return &SQLiteOGStmt{
		c:        c,
		sql:      query,
		numInput: strings.Count(query, "?"),
	}, nil
}

func (c *SQLiteOGConn) Close() error {
	_, err := c.OGClient.Close(context.Background(), &pb.ConnectionId{Id: c.ID})
	return err
}

func (c *SQLiteOGConn) Begin() (driver.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (c *SQLiteOGConn) Ping(ctx context.Context) error {
	_, err := c.OGClient.Ping(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	return nil
}

func (c *SQLiteOGConn) ResetSession(ctx context.Context) error {
	_, err := c.OGClient.ResetSession(ctx, &pb.ConnectionId{Id: c.ID})
	if err != nil {
		return err
	}
	return nil
}

func (c *SQLiteOGConn) IsValid() bool {
	_, err := c.OGClient.IsValid(context.Background(), &pb.ConnectionId{Id: c.ID})
	if err != nil {
		return false
	}
	return true
}

func (c *SQLiteOGConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	params, err := namedValuesToParams(args)
	if err != nil {
		return nil, err
	}
	stmt := &pb.Statement{
		Sql:    query,
		Params: params,
		CnxId:  c.ID,
	}
	pbr, err := c.OGClient.Execute(ctx, stmt)
	if err != nil {
		return nil, err
	}
	return resultFromPB(pbr)
}

func (c *SQLiteOGConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	params, err := namedValuesToParams(args)
	if err != nil {
		return nil, err
	}
	stmt := &pb.Statement{
		Sql:    query,
		Params: params,
		CnxId:  c.ID,
	}
	pbr, err := c.OGClient.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	return rowsFromPB(pbr)
}

func namedValuesToParams(namedValues []driver.NamedValue) ([]string, error) {
	// TODO: ignoring the "Name" now, this should be reviewed (see definition of driver.NamedValue)
	params := make([]string, len(namedValues))
	for _, v := range namedValues {
		if v.Ordinal < 1 {
			return nil, fmt.Errorf("ordinal cannot be < 1 %s %v", v.Name, v.Value)
		}
		params[v.Ordinal-1] = fmt.Sprintf("%v", v.Value)
	}
	return params, nil
}
