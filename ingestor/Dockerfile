FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ingestor ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/ingestor ./ingestor
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 50051

CMD ["./ingestor"]
