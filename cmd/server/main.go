package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/aousomran/sqlite-og/gen/proto"
	"github.com/aousomran/sqlite-og/internal/connections"
	"github.com/aousomran/sqlite-og/internal/server"
)

var (
	port             = flag.Int("port", 9091, "The server port")
	statsInterval    = flag.Duration("stats-interval", 0*time.Second, "interval in seconds for logging stats, use 0s to disable (default 0s)")
	logLevel         = flag.String("log-level", "info", "minimum log level to print out, choices (debug,info,warn,error) ")
	logFormat        = flag.String("log-format", "text", "log format choices (text,json)")
	pprofEnabled     = flag.Bool("enable-pprof", false, "enabled pprof at localhost:6060")
	discoveryEnabled = flag.Bool("enable-discovery", false, "enables grpc service discovery")
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	//log.Printf("\n\nRequest: %+v \n\n", req)
	slog.InfoContext(ctx, "request received", "method", info.FullMethod)
	defer func() {
		//log.Printf("\n\nResponse: %+v \n\n", resp)
		slog.InfoContext(ctx, "sent response")
		if err != nil {
			slog.ErrorCtx(ctx, "error", "reason", err)
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

func connectionStats(manager *connections.Manager, duration time.Duration) {
	// do not log connection stats
	if duration < 1*time.Second {
		slog.Info("connection stats disabled")
		return
	}
	ticker := time.NewTicker(duration).C
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

func initLogger(level, format string) {
	l := slog.Default()
	loggerOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	switch strings.ToLower(level) {
	case "debug":
		loggerOpts.Level = slog.LevelDebug
	case "warn":
		loggerOpts.Level = slog.LevelWarn
	case "error":
		loggerOpts.Level = slog.LevelError
	default:
		log.Printf("unknown log level %s, using default level info", level)
	}
	if format == "json" {
		l = slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts))
	} else {
		l = slog.New(slog.NewTextHandler(os.Stdout, loggerOpts))
	}
	slog.SetDefault(l)
}

func main() {
	flag.Parse()
	initLogger(*logLevel, *logFormat)

	// profiler
	if *pprofEnabled {
		go func() {
			if errProfiler := http.ListenAndServe("localhost:6060", nil); errProfiler != nil {
				slog.Error("could not start profiler", "error", errProfiler)
			}
		}()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	if *discoveryEnabled {
		reflection.Register(s)
	}

	manager := connections.NewManager()
	go connectionStats(manager, *statsInterval)
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

}
