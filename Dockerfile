FROM golang:1.19-alpine

WORKDIR /tinynamer

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN apk add zip make protobuf gcc musl-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go mod download && go mod verify
RUN export PATH="$PATH:$(go env GOPATH)/bin"
RUN echo $(ls $(go env GOPATH)/bin)

COPY . .
RUN make clean protobuf build

CMD ["./build/namerd"]
