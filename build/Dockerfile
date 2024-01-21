FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod .
COPY ./cmd/task-gopher/*.go ./
COPY ./data/* ./data/

RUN go get
RUN go build -o bin ./...

FROM build-stage AS run-tests-stage
RUN go test -v ./...

FROM run-tests-stage As deploy-stage

ENTRYPOINT [ "/app/bin", "serve"]
