package grpc2grpc

import (
	"context"
	"errors"
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
	petMapArr  map[int32]Pet
	userMapArr map[string]User
	Done       chan bool
}

var recUserArrClient = make([]User, 0)
var recUserArrServer = make([]User, 0)
var clntSendCount []int64
var clntRcvdCount [][]int64

//var srvrSendCount, srvrRcvdCount []int
//var srvrSendCount, srvrRcvdCount int
var clntEOFCount []int

//var srvrEOFCount int
var exitSignal bool

//var stopsignal bool

/* var mutex2 = &sync.Mutex{}
var mutex3 = &sync.Mutex{} */

var petArr = []Pet{
	{
		Id:   2,
		Name: "cat2",
	},
	{
		Id:   3,
		Name: "cat3",
	},
	{
		Id:   4,
		Name: "cat4",
	},
}
var userArr = []User{}

var sendUserDataClient = []User{}

var sendUserDataServer = []User{}

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

	// switch *option {
	// case "pet":
	// 	id, _ := strconv.Atoi(name)
	// 	return petById(client, id)
	// case "user":
	// 	return userByName(client, name)
	// case "listusers":
	// 	err := listUsers(client, data)
	// 	return nil, err
	// case "storeusers":
	// 	err := storeUsers(client, data)
	// 	return nil, err
	// case "bulkusers":

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
			//fmt.Println("for count: ", count, "threads", threads)
			if clntEOFCount[i] == 1 {
				count++
			} else {
				//fmt.Println("else count: ", count, "threads", threads)
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

				totalRspTm = totalRspTm + clntRcvdCount[i][0]*clntRcvdCount[i][1]
				if totalRspTm1 > totalRspTm {
					fmt.Println("Error#######################")
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

//userByName comments
func userByName(client PetStoreServiceClient, name string) (interface{}, error) {
	return nil, nil
}

//petById comments
func petById(client PetStoreServiceClient, id int) (interface{}, error) {

	return nil, nil
}

func listUsers(client PetStoreServiceClient, data func(data interface{}) bool) error {
	return nil
}

func storeUsers(client PetStoreServiceClient, data func(data interface{}) bool) error {
	return nil
}

func bulkUsers(client PetStoreServiceClient, data func(data interface{}) bool, thread int) error {
	stream, err := client.BulkUsers(context.Background())
	if err != nil {
		fmt.Println("err1", err)
		return err
	}
	var crc, rspTm, avg int64 = 0, 0, 0
	//var clntRspTm [][2]int64
	var csc int64 = 0

	waitc := make(chan struct{})
	go func() {
		for {
			//time.Sleep(time.Second)
			user, err := stream.Recv()
			if err == io.EOF {
				// read done.
				//fmt.Println("received total: ", strconv.Itoa(crc))

				/* for _, v := range clntRspTm {
					avg = avg + v[1]
				} */
				clntRcvdCount[thread][0] = crc
				clntRcvdCount[thread][1] = avg / crc
				clntEOFCount[thread] = 1
				recUserArrClient = nil
				close(waitc)
				//fmt.Println("CRC=", crc)
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
				rspTm = t1 - user.GetTimestmp()
				//fmt.Println("RECEIVED:  ", t1, user.GetTimestmp(), rspTm)
				if avg > avg+rspTm {
					fmt.Println("Error############=", avg, rspTm)
				}
				avg = avg + rspTm
				//clntRspTm = append(clntRspTm, [2]int64{crc, rspTm})
				//clntRcvdCount[thread][1] = (rspTm + ((crc - 1) * clntRcvdCount[thread][1])) / crc
				//fmt.Println("RECEIVED:  ", crc, rspTm, clntRcvdCount[thread][1])
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
			Timestmp: tN,
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
		petMapArr:  make(map[int32]Pet),
		userMapArr: make(map[string]User),
	}
	for _, pet := range petArr {
		server.petMapArr[pet.GetId()] = pet
	}

	for _, user := range userArr {
		server.userMapArr[""] = user
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

func (t *ServerStrct) PetById(ctx context.Context, req *PetByIdRequest) (*PetResponse, error) {
	return nil, errors.New("Pet not found")

}

func (t *ServerStrct) UserByName(ctx context.Context, req *UserByNameRequest) (*UserResponse, error) {

	return nil, errors.New("User not found")
}

func (t *ServerStrct) ListUsers(req *EmptyReq, sReq PetStoreService_ListUsersServer) error {
	return nil
}
func (t *ServerStrct) StoreUsers(cReq PetStoreService_StoreUsersServer) error {
	return nil
}
func (t *ServerStrct) BulkUsers(bReq PetStoreService_BulkUsersServer) error {
	fmt.Println("server BulkUsers method called")
	//var rc = 0
	//var st = 0
	/* waitr := make(chan struct{})
	go func() {
		for {
			userData, err := bReq.Recv()
			if userData != nil {
				rc++
				//fmt.Println("RECEIVED: ", userData.GetUsername())
				if t.Done != nil {
					t.Done <- true
				}
			}
			if err == io.EOF {
				 mutex.Lock()
				srvrEOFCount++
				srvrRcvdCount = append(srvrRcvdCount, rc)
				mutex.Unlock()
				fmt.Println("total records received: ", strconv.Itoa(rc))
				fmt.Println("server EOF count ", srvrEOFCount)

				mutex2.Lock()
				serverRec(rc)
				mutex2.Unlock()
				fmt.Println("total records received: ", strconv.Itoa(rc))
				fmt.Println("server EOF count ", srvrEOFCount)
				if srvrEOFCount == 10000 {
					fmt.Println("stopsignal: ", stopsignal)
					stopsignal = true
					fmt.Println("stopsignal: ", stopsignal)

				}
				close(waitr)
				return
			}
			if err != nil {
				fmt.Println("error occured while receiving", err)
				serverRec(rc)
				mutex1.Lock()
				srvrRcvdCount = append(srvrRcvdCount, rc)
				mutex1.Unlock()
				close(waitr)
				return
			}
		}
	}() */

	//fmt.Println("After Waitr", waitr)

	for {
		userData, err := bReq.Recv()
		if userData != nil {
			//rc++
			//fmt.Println("RECEIVED: ", userData.GetUsername())
			if t.Done != nil {
				t.Done <- true
			}
		}
		if err == io.EOF {
			/* mutex.Lock()
			srvrEOFCount++
			srvrRcvdCount = append(srvrRcvdCount, rc)
			mutex.Unlock() */
			/* fmt.Println("total records received: ", strconv.Itoa(rc))
			fmt.Println("server EOF count ", srvrEOFCount) */
			break
		}
		if err != nil {
			fmt.Println("error occured while receiving", err)
			break
		}
		/* user := User{
			Username: "suser1",
		} */
		//st++
		//fmt.Println("SEND: ", user.Username, st)
		if err := bReq.Send(userData); err != nil {
			fmt.Println("error while sending user", userData)
			//serverSent(st)
			return err
		}
	}
	//<-waitr
	//fmt.Println("srvrEOFCount waits", srvrEOFCount)
	if t.Done != nil {
		close(t.Done)
	}
	//time.Sleep(20 * time.Second)
	//fmt.Println("total client rcvd count: ", totalCount(srvrRcvdCount)
	//fmt.Println("total client send count: ", totalCount(srvrSendCount))
	/* fmt.Println("total client rcvd count: ", srvrRcvdCount)
	fmt.Println("total client send count: ", srvrSendCount) */
	return nil
}

/* func serverSent(n int) {
	//var mutex = &sync.Mutex{}
	//mutex.Lock()
	//time.Sleep(1 * time.Second)
	srvrSendCount = srvrSendCount + n
	//fmt.Println("grpcCallServerSent = \n", grpcCallsServerSent, time.Now())
	//mutex.Unlock()
}
func serverRec(n int) {
	//var mutex1 = &sync.Mutex{}
	//mutex1.Lock()
	//time.Sleep(1 * time.Second)
	fmt.Println("1srvrEOFCount =  srvrRcvdCount TIME: \n", srvrEOFCount, srvrRcvdCount, time.Now())
	srvrEOFCount++
	srvrRcvdCount = srvrRcvdCount + n
	fmt.Println("2srvrEOFCount =  srvrRcvdCount TIME: \n", srvrEOFCount, srvrRcvdCount, time.Now())
	//fmt.Println("grpcCallsServerRec =  TIME: \n", grpcCallsServerRec, time.Now())
	//mutex1.Unlock()
} */
