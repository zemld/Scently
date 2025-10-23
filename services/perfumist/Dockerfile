FROM golang:1.25-alpine

WORKDIR /perfumist

COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/main.go"]