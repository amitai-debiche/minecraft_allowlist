FROM golang:1.23 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/src/app/op_gateway op_gateway.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /usr/src/app/op_gateway .

EXPOSE 8080

CMD ["./op_gateway"]
