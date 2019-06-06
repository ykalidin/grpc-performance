package main

import (
	"flag"

	"github.com/project-flogo/core/engine"
	"github.com/ykalidin/grpc-performance/examples"
	"github.com/ykalidin/grpc-performance/proto/grpc2grpc"
)

func main() {
	clientMode := flag.Bool("client", false, "command to run client")
	serverMode := flag.Bool("server", false, "command to run server")
	number := flag.Int("number", 1, "number of times to run client")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")
	flag.Parse()

	if *clientMode {
		grpc2grpc.CallClient(port, *number, *paramVal)
		return
	} else if *serverMode {
		server, _ := grpc2grpc.CallServer()
		server.Start()
		return
	} else {

		e, err := examples.GRPC2GRPCExample()
		if err != nil {
			panic(err)
		}
		engine.RunEngine(e)
	}
}
