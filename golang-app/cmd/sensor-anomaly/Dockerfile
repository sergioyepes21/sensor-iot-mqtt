FROM golang:1.21 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o sensor-anomaly ./cmd/sensor-anomaly

CMD ["./sensor-anomaly"]