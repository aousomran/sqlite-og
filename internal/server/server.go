package server

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/internal/connections"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
)

type Server struct {
	pb.UnimplementedSqliteOGServer
	Manager *connections.Manager
}

func New(manager *connections.Manager) *Server {
	return &Server{
		Manager: manager,
	}
}

func toInterfaceSlice(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for k, v := range s {
		res[k] = v
	}
	return res
}

func (s *Server) Connection(ctx context.Context, in *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	id, err := s.Manager.Connect(in.GetDbName(), in.GetFunctions(), in.GetAggregators())
	if err != nil {
		return nil, err
	}
	return &pb.ConnectionResponse{
		ConnectionId: id,
	}, nil
}

func (s *Server) ExecuteQuery(ctx context.Context, in *pb.Statement) (*pb.ExecuteQueryResult, error) {
	if in.CnxId == "" {
		return nil, fmt.Errorf("connection id cannot be empty")
	}
	db, err := s.Manager.GetConnection(in.GetCnxId())
	if err != nil {
		slog.Error("cannot get database from manager", "error", err)
		return nil, err
	}
	fmt.Printf("params: %+v", in.GetParams())
	params := toInterfaceSlice(in.GetParams())
	columns, rows, err := db.Query(ctx, in.GetSql(), params...)
	if err != nil {
		return nil, err
	}

	return &pb.ExecuteQueryResult{
		Columns: columns,
		Rows:    rows,
		// we can't get those from DBWrapper.Query() it must be a DBWrapper.Execute() to work
		// this is a temporary measure ... don't want to use a regex as it's kinda slow
		LastInsertId: 0,
		AffectedRows: int64(len(rows)),
	}, nil
}

func (s *Server) CreateSQLiteFunction(ctx context.Context, in *pb.Signature) (*pb.CreateFunctionResult, error) {
	return &pb.CreateFunctionResult{Ok: true}, nil
}

func (s *Server) Callback(cbs pb.SqliteOG_CallbackServer) error {
	ctx := cbs.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		mdErr := errors.New("cannot get metadata from incoming context")
		slog.Error(mdErr.Error())
		return mdErr
	}
	cnxIdSlice, _ := md["cnx_id"]
	if len(cnxIdSlice) != 1 {
		mdErr := errors.New("connection id mismatch")
		slog.Error(mdErr.Error(), "want", 1, "got", len(cnxIdSlice))
		return mdErr
	}
	slog.Info("got connection id", "cnx_id", cnxIdSlice[0])
	db, err := s.Manager.GetConnection(cnxIdSlice[0])
	if err != nil {
		slog.Error("cannot get database from manager", "error", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				slog.Info("go routine 1 received done")
				wg.Done()
				break OUTER
			case invoke := <-db.Channels.ChanSend:
				slog.Info("sending invoke", "func_name", invoke.GetFunctionName(), "args", invoke.Args)
				errSend := cbs.Send(invoke)
				if errSend != nil {
					slog.Error("error sending invocation", "error", errSend.Error())
					wg.Done()
					break OUTER
				}
				continue
			}
		}
		slog.Info("server: end of go routine 1")
	}()

	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				slog.InfoCtx(ctx, "go routine 2 received done")
				wg.Done()
				break OUTER
			default:
				slog.Info("server: receiving")
				invokeResult, errRecv := cbs.Recv()
				slog.Info("received invoke result", "result", invokeResult.GetResult())
				if errRecv != nil {
					if errRecv == io.EOF {
						slog.Info("client finished streaming")
					} else {
						slog.Error("error receiving invocation result", "error", errRecv)
					}
					wg.Done()
					break OUTER
				}
				db.Channels.ChanReceive <- invokeResult
			}
		}
		slog.Info("server: end of go routine 2")
	}()

	wg.Wait()
	slog.Info("server: bye!")
	return nil
}

func (s *Server) Close(ctx context.Context, in *pb.CloseRequest) (*pb.Empty, error) {
	cnxId := in.GetConnectionId()
	if cnxId == "" {
		return nil, fmt.Errorf("got empty database connection")
	}
	connection, err := s.Manager.GetConnection(cnxId)
	if err != nil {
		return nil, fmt.Errorf("connection %s not found", cnxId)
	}
	errClose := connection.Close()
	if errClose != nil {
		return nil, fmt.Errorf("unable to close connection %d, error: %s", cnxId, errClose.Error())
	}
	s.Manager.DeleteConnection(cnxId)
	return &pb.Empty{}, nil
}

//func (s *Server) Query(ctx context.Context, in *pb.Statement) (*pb.QueryResult, error) {
//	fmt.Printf("params: %+v", in.GetParams())
//	params := toInterfaceSlice(in.GetParams())
//	rows, err := s.DBWrapper.QueryContext(ctx, in.GetSql(), params...)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	cols, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//
//	result := make([]*pb.Row, 0)
//	for rows.Next() {
//		r, err := rowStringSlice(cols, rows)
//		if err != nil {
//			return nil, err
//		}
//		result = append(result, &pb.Row{
//			Fields: r,
//		})
//	}
//
//	return &pb.QueryResult{
//		Success: true,
//		Columns: cols,
//		Rows:    result,
//	}, nil
//}
//
//func (s *Server) Execute(ctx context.Context, in *pb.Statement) (*pb.ExecuteResult, error) {
//	stmt, err := s.DBWrapper.PrepareContext(ctx, in.Sql)
//	if err != nil {
//		return nil, err
//	}
//	defer stmt.Close()
//
//	result, err := stmt.ExecContext(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	insertId, err := result.LastInsertId()
//	if err != nil {
//		return nil, err
//	}
//	affected, err := result.RowsAffected()
//	if err != nil {
//		return nil, err
//	}
//
//	return &pb.ExecuteResult{
//		Success:      true,
//		LastInsertId: insertId,
//		AffectedRows: affected,
//	}, nil
//}
