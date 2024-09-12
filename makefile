# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Default number of clusters
Njob ?= 2
Ncomet ?= 2

# Build all jobs
all: test build

# Build job directories and binaries
build-job:
	mkdir -p bin/job; \
	rm -rf bin/job/job; \
	$(GOBUILD) -o bin/job/job ./src/cmd/job/main.go; \

# Run all jobs
run-job:
	for i in $(shell seq 1 $(Njob)); do \
		echo "Running job$$i..."; \
		(cd bin/job && ./job -c ../../config/job$$i/config.yml) & \
	done

stop-job:
	pkill -f ./job



build-comet:
	mkdir -p bin/comet; \
	rm -rf bin/comet/comet; \
	$(GOBUILD) -o bin/comet/comet ./src/cmd/comet/main.go; \

# Run all comets
run-comet:
	for i in $(shell seq 1 $(Ncomet)); do \
		echo "Running comet$$i..."; \
		(cd bin/comet && ./comet -c ../../config/comet$$i/config.yml) & \
	done

stop-comet:
	pkill -f ./comet


test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

server:
	 goreman -f local-cluster-profile start; \
	 bash redisClusterStart.sh; \
	 /opt/kafka/bin/kafka-server-start.sh ./config/kafka.properties; 

proto:
	protoc --go_out=.. --go-grpc_out=.. --go-grpc_opt=require_unimplemented_servers=false -I. -Iproto proto/msg/msg.proto proto/comet/comet.proto proto/logic/logic.proto; 
# run:
# 	nohup target/logic -conf=target/logic.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
# 	nohup target/comet -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &
# 	nohup target/job -conf=target/job.toml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/job.log &

# stop:
# 	pkill -f target/logic
# 	pkill -f target/job
# 	pkill -f target/comet
