FROM golang:1.25.1-alpine AS builder

ENV PATH=$PATH:/root/go/bin

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o perfume-hub ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/perfume-hub .

EXPOSE 8000

CMD ["./perfume-hub"]
