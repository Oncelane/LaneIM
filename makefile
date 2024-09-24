# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

# Default number of clusters
Njob ?= 1
Ncomet ?= 2

Njobp ?= 1
Ncometp ?= 1

Nlogic ?= 1
Nlogicp ?= 1

# Build all jobs
all: stop-job stop-comet build-job build-comet

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
run-jobp:
	for i in $(shell seq 1 $(Njobp)); do \
		echo "Running job$$i..."; \
		(cd bin/job && ./job -c ../../pConfig/job$$i/config.yml) & \
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
run-cometp:
	for i in $(shell seq 1 $(Ncometp)); do \
		echo "Running comet$$i..."; \
		(cd bin/comet && ./comet -c ../../pConfig/comet$$i/config.yml) & \
	done

stop-comet:
	pkill -f ./comet

build-logic:
	mkdir -p bin/logic; \
	rm -rf bin/logic/logic; \
	$(GOBUILD) -o bin/logic/logic ./src/cmd/logic/main.go; \

# Run all logics
run-logic:
	for i in $(shell seq 1 $(Nlogic)); do \
		echo "Running logic$$i..."; \
		(cd bin/logic && ./logic -c ../../config/logic$$i/config.yml) & \
	done
run-logicp:
	for i in $(shell seq 1 $(Nlogicp)); do \
		echo "Running logic$$i..."; \
		(cd bin/logic && ./logic -c ../../pConfig/logic$$i/config.yml) & \
	done

stop-logic:
	pkill -f ./logic
	
test:
	$(GOTEST) -v ./...

clean:
	rm -rf bin/
# run:
# 	nohup target/logic -conf=target/logic.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
# 	nohup target/comet -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &
# 	nohup target/job -conf=target/job.toml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/job.log &

# stop:
# 	pkill -f target/logic
# 	pkill -f target/job
# 	pkill -f target/comet
