package main

import (
	"context"
	"flag"
	"fmt"
	pb "gitlab.com/aous.omr/sqlite-og/gen/proto"
	"gitlab.com/aous.omr/sqlite-og/internal/connections"
	"gitlab.com/aous.omr/sqlite-og/internal/server"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("\n\nRequest: %+v \n\n", req)
	defer func() {
		log.Printf("\n\nResponse: %+v \n\n", resp)
		if err != nil {
			log.Printf("\n\nError: %+v \n\n", err)
		}
	}()

	return handler(ctx, req)
}

//func handleSignals(closers ...io.Closer) {
//	c := make(chan os.Signal)
//	signal.Notify(c)
//	select {
//		case sig := <-c:
//			log.Printf("os signal %+v recieved", sig)
//			for _,v := range closers {
//				_ := v.close()
//			}
//	}
//}

func connectionStats(manager *connections.Manager) {
	ticker := time.NewTicker(2 * time.Second).C
	for {
		select {
		case <-ticker:
			keys := make([]string, 0)
			for k, _ := range manager.CnxMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			slog.Info("db connection stats", "count", len(keys), "ids", keys)
		}
	}
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	reflection.Register(s)

	manager := connections.NewManager()
	go connectionStats(manager)
	srv := server.New(manager)
	pb.RegisterSqliteOGServer(s, srv)
	slog.Info("server listening ", "addr", listener.Addr())

	// close things gracefully
	defer func() {
		s.GracefulStop()
		err := manager.Close()
		if err != nil {
			slog.Warn("at least one connection failed to close", "error", err.Error())
		}

		err = listener.Close()
		if err != nil {
			slog.Warn("failed to close listener", "error", err.Error())
		}
	}()

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	slog.Info("server: bye!")

}
