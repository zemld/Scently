FROM golang:1.25-alpine AS builder

RUN apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main protoc protobuf-dev make || \
    apk add --no-cache protobuf protobuf-dev make

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH=$PATH:/root/go/bin

WORKDIR /app

COPY schemas/ ./schemas/
COPY Makefile ./

RUN make gen-perfume-hub
RUN make gen-perfumist

COPY generated/go.mod generated/go.sum ./generated/

COPY services/perfumist/go.mod services/perfumist/go.sum ./services/perfumist/
WORKDIR /app/services/perfumist
RUN go mod download

COPY services/perfumist/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o perfumist ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/perfumist/perfumist .

EXPOSE 8000

CMD ["./perfumist"]