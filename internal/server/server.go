package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc/metadata"
	"vitess.io/vitess/go/vt/sqlparser"

	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/internal/connections"
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

func (s *Server) Connection(ctx context.Context, in *pb.ConnectionRequest) (*pb.ConnectionId, error) {
	id, err := s.Manager.Connect(in.GetDbName(), in.GetFunctions(), in.GetAggregators())
	if err != nil {
		return nil, err
	}
	return &pb.ConnectionId{
		Id: id,
	}, nil
}

func (s *Server) ExecuteOrQuery(ctx context.Context, in *pb.Statement) (*pb.ExecuteOrQueryResult, error) {
	if in.CnxId == "" {
		return nil, fmt.Errorf("connection id cannot be empty")
	}

	result := &pb.ExecuteOrQueryResult{
		QueryResult: &pb.QueryResult{
			Columns:     []string{},
			ColumnTypes: []string{},
			Rows:        []*pb.Row{},
		},
		ExecuteResult: &pb.ExecuteResult{
			LastInsertId: -1,
			AffectedRows: -1,
		},
	}

	// this is a DML operation, we'll call execute & fill ExecuteResult
	if sqlparser.IsDML(in.GetSql()) {
		slog.InfoContext(ctx, "executing", "isDML", true, "sql", fmt.Sprintf("%s", in.GetSql()), "params", fmt.Sprintf("%+v", in.GetParams()))
		execute, err := s.Execute(ctx, in)
		if err != nil {
			return nil, err
		}
		result.ExecuteResult.LastInsertId = execute.GetLastInsertId()
		result.ExecuteResult.AffectedRows = execute.GetAffectedRows()

	} else {
		slog.InfoContext(ctx, "querying", "isDML", false, "sql", fmt.Sprintf("%s", in.GetSql()), "params", fmt.Sprintf("%+v", in.GetParams()))
		query, err := s.Query(ctx, in)
		if err != nil {
			return nil, err
		}

		result.QueryResult.Columns = query.GetColumns()
		result.QueryResult.ColumnTypes = query.GetColumnTypes()
		result.QueryResult.Rows = query.GetRows()
		// temporary .. to see if this fixes django
		result.ExecuteResult.AffectedRows = int64(len(query.Rows))
	}

	return result, nil
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
	slog.InfoContext(ctx, "got connection id", "cnx_id", cnxIdSlice[0])
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
				slog.Debug("callback go routine 1 received done")
				wg.Done()
				break OUTER
			case invoke := <-db.Channels.ChanSend:
				slog.Debug("sending invoke", "func_name", invoke.GetFunctionName(), "args", invoke.Args)
				errSend := cbs.Send(invoke)
				if errSend != nil {
					slog.Error("error sending invocation", "error", errSend.Error())
					wg.Done()
					break OUTER
				}
				continue
			}
		}
		slog.Debug("callback end of go routine 1")
	}()

	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				slog.Debug("go routine 2 received done")
				wg.Done()
				break OUTER
			default:
				slog.Debug("server: receiving")
				invokeResult, errRecv := cbs.Recv()
				slog.Debug("received invoke result", "result", invokeResult.GetResult())
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
		slog.Debug("server: end of go routine 2")
	}()

	wg.Wait()
	return nil
}

func (s *Server) Close(ctx context.Context, in *pb.ConnectionId) (*pb.Empty, error) {
	cnxId := in.GetId()
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

func (s *Server) IsValid(ctx context.Context, in *pb.ConnectionId) (*pb.Empty, error) {
	_, err := s.Manager.GetConnection(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	// All we want for now is that this server is reachable
	return &pb.Empty{}, nil
}

func (s *Server) ResetSession(ctx context.Context, in *pb.ConnectionId) (*pb.ConnectionId, error) {
	// NoOp for now, let's just make sure that the connection is valid
	_, err := s.IsValid(ctx, in)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (s *Server) Query(ctx context.Context, in *pb.Statement) (*pb.QueryResult, error) {
	db, err := s.Manager.GetConnection(in.GetCnxId())
	if err != nil {
		slog.Error("cannot get database from manager", "error", err)
		return nil, err
	}

	params := toInterfaceSlice(in.GetParams())
	columns, columnTypes, rows, err := db.Query(ctx, in.GetSql(), params...)
	if err != nil {
		slog.Error("error calling db.Query", "error", err.Error())
		return nil, err
	}

	return &pb.QueryResult{
		Columns:     columns,
		ColumnTypes: columnTypes,
		Rows:        rows,
	}, nil
}

func (s *Server) Execute(ctx context.Context, in *pb.Statement) (*pb.ExecuteResult, error) {
	db, err := s.Manager.GetConnection(in.GetCnxId())
	if err != nil {
		slog.Error("cannot get database from manager", "error", err)
		return nil, err
	}
	params := toInterfaceSlice(in.GetParams())

	lastInsertId, affectedRows, err := db.Execute(ctx, in.GetSql(), params...)
	if err != nil {
		slog.Error("error calling db.Execute", "error", err.Error())
		return nil, err
	}

	return &pb.ExecuteResult{
		LastInsertId: lastInsertId,
		AffectedRows: affectedRows,
	}, nil
}
