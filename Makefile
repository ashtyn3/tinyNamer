all: clean protobuf build run 

clean:
	rm -r build

protobuf:
	protoc protobufs/*.proto --go_out=./protobufs

build:
	mkdir build
	go build -o build/namerd ./cmd/main.go

run:
	./build/namerd

docker:
	docker build . --tag tinynamer

