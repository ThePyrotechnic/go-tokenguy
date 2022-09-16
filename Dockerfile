FROM golang:alpine3.16

WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY main.go .
COPY tokenguy/ ./tokenguy

RUN ls
RUN go mod download
RUN go build -o ./go-tokenguy

FROM alpine:3.16

WORKDIR /app
COPY --from=0 /app/go-tokenguy .

COPY keys ./keys
COPY config ./config

ENV GIN_MODE=release

EXPOSE 6666

ENTRYPOINT ["./go-tokenguy", "server", "start"]
