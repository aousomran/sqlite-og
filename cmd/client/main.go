package main

import (
	"context"
	"fmt"
	pb "github.com/aousomran/sqlite-og/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
)

func djangoDatetimeTrunc(args ...string) string {
	if len(args) < 2 {
		return ""
	}
	return strings.Split(args[1], " ")[0]
}

func square(args ...string) string {
	if len(args) < 1 {
		return ""
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("cannot convert %s to integer", args[0])
		return ""
	}
	return strconv.Itoa(i)
}

type callBackFunc func(args ...string) string

var (
	functionRegister = map[string]callBackFunc{
		"django_datetime_trunc": djangoDatetimeTrunc,
		"square":                square,
	}
)

func query(ctx context.Context, client pb.SqliteOGClient, cnxId string, sql string, params []string) {
	log.Println("starting query")
	Res, err := client.ExecuteQuery(ctx, &pb.Statement{
		CnxId:  cnxId,
		Sql:    sql,
		Params: params,
	})
	if err != nil {
		log.Fatalf("cannot execute query %v", err)
	}
	log.Println(strings.Join(Res.Columns, "\t\t"))
	log.Println(strings.Repeat("--", 40))
	for _, v := range Res.Rows {
		log.Println(fmt.Sprintf(strings.Join(v.Fields, "\t\t")))
	}
}

func doCallbackDance(ctx context.Context, client pb.SqliteOGClient) *sync.WaitGroup {
	callbackClient, err := client.Callback(ctx)
	if err != nil {
		log.Fatalf("cannot create callback client %s\n", err.Error())
	}

	receiveChannel := make(chan *pb.Invoke)
	wg := sync.WaitGroup{}
	var msgsSent, msgsReceived int

	wg.Add(2)
	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				log.Println("1st go routine received ctx.Done()")
				wg.Done()
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
			case <-ctx.Done():
				log.Println("2nd go routine received ctx.Done()")
				break OUTER
			case invoke := <-receiveChannel:
				funcName := invoke.GetFunctionName()
				funk, ok := functionRegister[funcName]
				if !ok {
					log.Fatalf("requested function name that does not exist %s", funcName)
				}
				evaluate := funk(invoke.Args...)
				errSend := callbackClient.Send(&pb.InvocationResult{
					Initial: false,
					Result:  []string{evaluate},
				})
				if errSend != nil {
					log.Printf("got an error sending invocation result %v\n", errSend)
				}
				msgsSent++
			}
		}
	}()

	return &wg
}

func main() {

	//ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)
	go func() {
		select {
		case sig := <-sigChan:
			log.Printf("recevied signal %+v", sig)
			//cancel()
			break
		}
	}()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot create grpc connection %v\n", err)
	}

	client := pb.NewSqliteOGClient(conn)

	functionNames := make([]string, 0)
	for k, _ := range functionRegister {
		functionNames = append(functionNames, k)
	}

	cnx, err := client.Connection(ctx, &pb.ConnectionRequest{
		Functions:   functionNames,
		Aggregators: nil,
	})
	if err != nil {
		log.Fatalf("cannot create database connection %s\n", err.Error())
	}

	log.Printf("got cnx id %s", cnx.GetConnectionId())
	ctx = metadata.AppendToOutgoingContext(ctx, "cnx_id", cnx.GetConnectionId())

	_ = doCallbackDance(ctx, client)

	sql := `SELECT django_datetime_trunc(?, "1990-01-01 20:50:55", ?, ?)`
	params := []string{"day", "UTC", "UTC"}
	//sql := `SELECT * from app_channel limit 10`
	//params := []string{}
	query(ctx, client, cnx.GetConnectionId(), sql, params)

	//sql := `SELECT DISTINCT django_datetime_trunc(?, "app_channelstats"."date", ?, ?) AS "scrape_date" FROM "app_channelstats" ORDER BY "scrape_date" DESC LIMIT 10`
	////sql := `SELECT id, django_datetime_trunc(?, "app_channelstats"."date", ?, ?) AS "scrape_date" FROM "app_channelstats" ORDER BY "scrape_date" DESC LIMIT 3`
	////sql := `SELECT django_datetime_trunc(?, "1990-01-01 20:50:55", ?, ?)`
	//log.Println(strings.Repeat("###", 20))
	//for i := 0; i < 100; i++ {
	//	query(ctx, client, sql, []string{"day", "UTC", "UTC"})
	//}
	////cancel()
	//errClose := conn.Close()
	//if errClose != nil {
	//	log.Printf("error closing channel %+v\n", errClose)
	//}
	//log.Println(strings.Repeat("###", 20))
	//log.Printf("Received: %d\t\t Sent: %d\n", msgsReceived, msgsSent)
	//log.Println("client: bye!")
}

//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//
//	sigChan := make(chan os.Signal)
//	signal.Notify(sigChan)
//	go func() {
//		for {
//			log.Println("waiting for signal")
//			select {
//			case sig := <-sigChan:
//				log.Printf("got signal %+v", sig)
//				cancel()
//				break
//			}
//			break
//		}
//	}()
//
//	wg := sync.WaitGroup{}
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			num := rand.Int()
//			select {
//			case <-ctx.Done():
//				log.Printf("done from routine %d", num)
//				wg.Done()
//				break
//			}
//
//		}()
//	}
//	wg.Wait()
//	log.Println("client bye!")
//}