all: clean protobuf build run 

clean:
	rm -r build
build:
	mkdir build
	go build -o build/namerd ./cmd/namerd/main.go 

run:
	./build/namerd

protobuf:
	protoc protobufs/peer.proto --go_out=./protobufs

