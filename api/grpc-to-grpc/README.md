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
1)Start sample gRPC server:
```bash
./grpc-to-grpc -server
```

2)Start proxy gateway:

Note: If the gRPC server is running in a different system please set the GRPCSERVERIP variable to corresponding system IP before starting the proxy gateway.
```bash
export GRPCSERVERIP=XX.XX.XX.XX
```

```bash
export FLOGO_LOG_LEVEL=ERROR
FLOGO_RUNNER_TYPE=DIRECT ./grpc-to-grpc
```

3)Run gRPC client:

Note: If the proxy gateway is running in a different system please set the GWIP variable to corresponding system IP before starting the client.

```bash
export GWIP=XX.XX.XX.XX
```

```bash
./grpc-to-grpc -client -port 9096 -param 3 -number 2000
``` 

The test will be completed in 2 minutes and the results will be displayed in the client terminal.
