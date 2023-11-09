package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/aousomran/sqlite-og/gen/proto"
)

func init() {
	sql.Register("sqliteog", &SQLiteOGDriver{})
}

type SQLiteOGConnector struct {
	driver *SQLiteOGDriver
	host   string
	port   string
	tls    bool
	dbname string
}

func (c *SQLiteOGConnector) Connect(ctx context.Context) (driver.Conn, error) {
	target := fmt.Sprintf("%s:%s", c.host, c.port)
	grpcConn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewSqliteOGClient(grpcConn)
	cnxId, err := client.Connection(ctx, &pb.ConnectionRequest{
		DbName:      c.dbname,
		Functions:   nil,
		Aggregators: nil,
	})
	if err != nil {
		return nil, err
	}
	return &SQLiteOGConn{
		ID:       cnxId.GetId(),
		GRPCConn: grpcConn,
		OGClient: client,
	}, nil
}

func (c *SQLiteOGConnector) Driver() driver.Driver {
	return c.driver
}

type SQLiteOGDriver struct{}

func (c *SQLiteOGDriver) Open(dsn string) (driver.Conn, error) {
	ctr, err := c.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}
	return ctr.Connect(context.Background())
}

func (c *SQLiteOGDriver) OpenConnector(dsn string) (driver.Connector, error) {
	s1 := strings.Split(dsn, "/")
	if len(s1) < 2 {
		return nil, fmt.Errorf("wrong dsn format, must be `host:port/dbname`, got `%s`", dsn)
	}
	s2 := strings.Split(s1[0], ":")
	if len(s2) < 2 {
		return nil, fmt.Errorf("wrong dsn format, must be `host:port/dbname`, got `%s`", dsn)
	}
	host, port, dbname := s2[0], s2[1], s1[1]
	return &SQLiteOGConnector{
		driver: c,
		host:   host,
		port:   port,
		dbname: dbname,
		tls:    false,
	}, nil
}
