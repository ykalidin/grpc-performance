# gRPC to gRPC
This recipe is a proxy gateway for gRPC end points.

## Installation
* Install [Go](https://golang.org/)

## Setup
```bash
git clone https://github.com/ykalidin/grpc-performance.git
cd grpc-performance/api/grpc-to-grpc
go build
```

## Testing
1)Start proxy gateway:
```bash
export FLOGO_LOG_LEVEL=ERROR
FLOGO_RUNNER_TYPE=DIRECT ./grpc-to-grpc
```

2)Start sample gRPC server in another terminal.
```bash
./grpc-to-grpc -server
```

3)Now run gRPC client in new terminal.
```bash
./grpc-to-grpc -client -port 9096 -param 3 -number 2000
``` 

The test will be completed in 2 minutes and the results will be displayed in the client terminal.
