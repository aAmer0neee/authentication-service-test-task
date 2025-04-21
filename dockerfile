FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod tidy
RUN go build -o server ./cmd/main.go

CMD ["./server", "--config-path=./config.yaml"]