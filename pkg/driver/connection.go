package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"

	pb "github.com/aousomran/sqlite-og/gen/proto"
)

type SQLiteOGConn struct {
	ID       string
	DBName   string
	GRPCConn *grpc.ClientConn
	OGClient pb.SqliteOGClient
	Funcs    map[string]callbackFunc
}

func NewConnection(ctx context.Context, dbname string, grpcConn *grpc.ClientConn, callbacksEnabled bool, callbacks map[string]callbackFunc) (*SQLiteOGConn, error) {
	client := pb.NewSqliteOGClient(grpcConn)
	funcs := make(map[string]callbackFunc)
	var funcNames []string
	if callbacksEnabled {
		for k, v := range callbacks {
			funcs[k] = v
			funcNames = append(funcNames, k)
		}
	}

	cnxId, err := client.Connection(ctx, &pb.ConnectionRequest{
		DbName:      dbname,
		Functions:   funcNames,
		Aggregators: nil,
	})

	if err != nil {
		return nil, err
	}

	cnx := &SQLiteOGConn{
		ID:       cnxId.Id,
		DBName:   dbname,
		GRPCConn: grpcConn,
		OGClient: client,
		Funcs:    funcs,
	}

	if callbacksEnabled {
		err = cnx.DoCallbackDance(ctx)
		if err != nil {
			// attempt to close the connection
			//_, _ = client.Close(ctx, cnxId)
			return nil, err
		}
	}

	return cnx, nil
}

func (c *SQLiteOGConn) DoCallbackDance(ctx context.Context) error {
	//ctx = metadata.AppendToOutgoingContext(ctx, "cnx_id", c.ID)
	ctx2 := context.Background()
	ctx2 = metadata.AppendToOutgoingContext(ctx, "cnx_id", c.ID)
	callbackClient, err := c.OGClient.Callback(ctx2)
	if err != nil {
		//log.Fatalf("cannot create callback client %s\n", err.Error())
		return err
	}

	receiveChannel := make(chan *pb.Invoke)
	//c.wg = &sync.WaitGroup{}

	var msgsSent, msgsReceived int

	//c.wg.Add(2)
	go func() {
	OUTER:
		for {
			select {
			case <-ctx2.Done():
				log.Println("1st go routine received ctx.Done()")
				//c.wg.Done()
				break OUTER
			default:
				invoke, invErr := callbackClient.Recv()
				if invErr != nil {
					log.Printf("error receiving %v\n", invErr)
				}
				receiveChannel <- invoke
				msgsReceived++
			}
		}
	}()

	go func() {
	OUTER:
		for {
			select {
			case <-ctx2.Done():
				log.Println("2nd go routine received ctx.Done()")
				break OUTER
			case invoke := <-receiveChannel:
				funcName := invoke.GetFunctionName()
				callable, ok := c.Funcs[funcName]
				if !ok {
					log.Fatalf("requested function name that does not exist %s", funcName)
				}
				evaluate := callable(invoke.Args...)
				errSend := callbackClient.Send(&pb.InvocationResult{
					Initial: false,
					Result:  evaluate,
				})
				if errSend != nil {
					log.Printf("got an error sending invocation result %v\n", errSend)
				}
				msgsSent++
			}
		}
	}()

	return nil
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
