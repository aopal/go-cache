.PHONY: build clean

run: build run-built
log: build run-log

build:
	go build -v -o build/ ./cmd/main.go

run-built:
	./build/main config/config.default.yaml

run-log:
	./build/main config/config.default.yaml > logfile 2>&1

clean:
	rm build/*
