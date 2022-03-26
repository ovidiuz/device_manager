FROM golang:1.18-alpine

RUN apk update && apk add git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service

CMD ["./service"]