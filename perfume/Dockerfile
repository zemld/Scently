FROM golang:1.25-alpine

WORKDIR /perfume

COPY . .

RUN go mod tidy

CMD ["go", "run", "core/main.go"]