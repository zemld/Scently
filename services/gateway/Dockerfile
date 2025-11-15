FROM golang:1.25.1-alpine

WORKDIR /gateway

COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/main.go"]