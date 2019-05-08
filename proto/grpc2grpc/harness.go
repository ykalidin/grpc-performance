package grpc2grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
	Server     *grpc.Server
	userMapArr map[string]User
}

var clntSendCount []int64
var clntRcvdCount [][]int64

var clntEOFCount []int

var exitSignal bool

func CallClient(port *string, option *string, name string, data func(data interface{}) bool, threads int) (interface{}, error) {
	clientAddr := *port
	if len(*port) == 0 {
		clientAddr = ":9000"
	}
	clientAddr = ":" + *port
	conn, err := grpc.Dial("localhost"+clientAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewPetStoreServiceClient(conn)

	clntSendCount = make([]int64, threads)
	clntRcvdCount = make([][]int64, threads)
	for i := range clntRcvdCount {
		clntRcvdCount[i] = make([]int64, 2)
		fmt.Println(i)
	}
	clntEOFCount = make([]int, threads)
	exitSignal = false

	fmt.Println("Starting threads", time.Now())

	for thread := 0; thread < threads; thread++ {
		time.Sleep(20 * time.Millisecond)
		fmt.Println(thread)
		go bulkUsers(client, data, thread)
	}
	fmt.Println("All Threads Activated", time.Now())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	fmt.Println("Exit signal received", time.Now())
	exitSignal = true
	var totalRspTm, totalRspTm1, clntSent, clntRcvd int64
	for {
		time.Sleep(time.Second)
		count := 1
		for i := 0; i < threads; i++ {
			if clntEOFCount[i] == 1 {
				count++
			} else {
				count = 1
				break
			}
			fmt.Println("count: ", count, "threads", threads)
			if count == threads {
				fmt.Println("All threads completed")
				break
			}
		}
		if count != 1 {
			fmt.Println("All threads completed")

			for i := 0; i < threads; i++ {

				totalRspTm = totalRspTm + clntRcvdCount[i][1]
				if totalRspTm1 > totalRspTm {
					fmt.Println("@@@@@@@@@@@@@ Overflow Error @@@@@@@@@@@@@", time.Now())
					os.Exit(1)
				}
				clntSent = clntSent + clntSendCount[i]
				clntRcvd = clntRcvd + clntRcvdCount[i][0]
				totalRspTm1 = totalRspTm
			}
			fmt.Println("totalGrpcCallRcv: ", clntRcvd)
			fmt.Println("totalGrpcCallSent: ", clntSent)
			fmt.Println("AverageResponseTime:", totalRspTm/clntRcvd)
			fmt.Println("Threads:", threads)
			fmt.Println("TimeBtwnMessages:30Ms")
			fmt.Println("TimeBtwnThread:10Ms")

			return nil, nil
		}
	}

}

func totalCount(arr []int) int {
	total := 0
	for i := 0; i < len(arr); i++ {
		total = total + arr[i]
	}
	return total
}

func bulkUsers(client PetStoreServiceClient, data func(data interface{}) bool, thread int) error {
	stream, err := client.BulkUsers(context.Background())
	if err != nil {
		fmt.Println("err1", err)
		return err
	}
	var crc, rspTm, ttlRspTm int64
	var csc int64

	waitc := make(chan struct{})
	go func() {
		for {
			user, err := stream.Recv()
			if err == io.EOF {
				clntRcvdCount[thread][0] = crc
				clntRcvdCount[thread][1] = ttlRspTm
				clntEOFCount[thread] = 1
				close(waitc)

				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a user : %v", err)
				close(waitc)
				return
			}
			if user != nil {
				crc++
				t1 := time.Now().UnixNano()
				rspTm = t1 - user.GetTimestamp1()
				if ttlRspTm > ttlRspTm+rspTm {
					fmt.Println("@@@@@@@@@@@@@ Overflow Error @@@@@@@@@@@@@", ttlRspTm, rspTm)
					os.Exit(1)
				}
				ttlRspTm = ttlRspTm + rspTm
				if data != nil {
					data(user)
				}
			}
		}
	}()
	for {
		tN := time.Now().UnixNano()
		user := User{
			// Username: "cuser" + strconv.Itoa(thread),
			Timestamp1: tN,
		}
		//fmt.Println("SEND: ", user.Username)
		if err := stream.Send(&user); err != nil {
			fmt.Println("error while sending user", user, err)
			return err
		}
		csc++
		time.Sleep(30 * time.Millisecond)
		if exitSignal {
			fmt.Println("Send Close signal recieved")
			clntSendCount[thread] = csc
			break
		}
	}
	fmt.Println("Closing for thread: ", thread)
	stream.CloseSend()
	<-waitc

	return nil
}

func CallServer() (*ServerStrct, error) {
	s := grpc.NewServer()
	server := ServerStrct{
		Server:     s,
		userMapArr: make(map[string]User),
	}

	RegisterPetStoreServiceServer(s, &server)

	return &server, nil
}

func (t *ServerStrct) Start() error {
	addr := ":9000"
	log.Println("Starting server on port: ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("err2", err)
		return err
	}

	return t.Server.Serve(lis)
}

func (t *ServerStrct) BulkUsers(bReq PetStoreService_BulkUsersServer) error {
	fmt.Println("server BulkUsers method called")
	for {
		userData, err := bReq.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error occured while receiving", err)
			break
		}

		if err := bReq.Send(userData); err != nil {
			fmt.Println("error while sending user", userData)
			return err
		}
	}
	return nil
}
