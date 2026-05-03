FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o go-vless-client .

FROM alpine:3.21

COPY --from=builder /app/go-vless-client /usr/local/bin/go-vless-client

ENTRYPOINT ["go-vless-client"]
