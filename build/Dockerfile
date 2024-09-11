FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o short-video-aggregator cmd/main.go

EXPOSE 8000

CMD ["./short-video-aggregator"]
