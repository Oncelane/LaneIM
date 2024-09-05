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
build-job-group:
	mkdir -p bin
	for i in $(shell seq 1 $(Njob)); do \
		rm -rf bin/job$$i; \
		mkdir -p bin/job$$i; \
		cp src/cmd/job/job.yml bin/job$$i/job.yml; \
		$(GOBUILD) -o bin/job$$i/job ./src/cmd/job/main.go; \
	done

# Run all jobs
run-job-group:
	for i in $(shell seq 1 $(Njob)); do \
		echo "Running job$$i..."; \
		(cd bin/job$$i && ./job) & \
	done

stop-job-group:
	pkill -f ./job



build-comet-group:
	mkdir -p bin
	for i in $(shell seq 1 $(Ncomet)); do \
		rm -rf bin/comet$$i; \
		mkdir -p bin/comet$$i; \
		cp src/cmd/comet/comet.yml bin/comet$$i/comet.yml; \
		$(GOBUILD) -o bin/comet$$i/comet ./src/cmd/comet/main.go; \
	done

# Run all comets
run-comet-group:
	for i in $(shell seq 1 $(Ncomet)); do \
		echo "Running comet$$i..."; \
		(cd bin/comet$$i && ./comet) & \
	done

stop-comet-group:
	pkill -f ./comet


test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

# run:
# 	nohup target/logic -conf=target/logic.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
# 	nohup target/comet -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &
# 	nohup target/job -conf=target/job.toml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/job.log &

stop:
	pkill -f target/logic
	pkill -f target/job
	pkill -f target/comet
