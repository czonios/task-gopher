FROM golang:1.21

WORKDIR /app

COPY go.mod .
COPY ./cmd/task-gopher/task-gopher.go .
COPY ./cmd/task-gopher/cli.go .
COPY ./cmd/task-gopher/server.go .

RUN go get
RUN go build -o bin ./...

ENTRYPOINT [ "/app/bin", "serve"]