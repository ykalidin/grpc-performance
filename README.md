# grpc-performance
This is a gRPC Bidirectional Streaming example.The messages stream from gRPC client-->Gateway-->gRPC Echo Server-->Gateway-->gRPC Client.

The gRPC Client outputs the Number of messages sent and received, Average Response Time taken by all the requests, Moving Average Response Time and Transactions per second values.

## Installation
* The Go programming language version 1.11 or later should be [installed](https://golang.org/doc/install).

## Setup
```
git clone https://github.com/ykalidin/grpc-performance
go get -u github.com/project-flogo/cli/...
go get -u github.com/project-flogo/microgateway/...
go get -u github.com/rameshpolishetti/movingavg/...
cd ykalidin/grpc-performance/api/grpc-to-grpc/
go build
```

## Testing
1)Start proxy gateway:
```bash
ulimit -n 1048576
export FLOGO_LOG_LEVEL=ERROR
FLOGO_RUNNER_TYPE=DIRECT ./grpc-to-grpc
```

2)Start sample gRPC server in a new terminal.
```bash
ulimit -n 1048576
./grpc-to-grpc -server
```

3)Start the gRPC client in a new terminal.
```bash
ulimit -n 1048576
./grpc-to-grpc -client -method bulkusers -port 9096 -number 1000
```

4)Stop the gRPC client after the required test duration.
All the outputs are printed in the client Terminal.
