FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main protoc protobuf-dev make || \
    apk add --no-cache protobuf protobuf-dev make

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH=$PATH:/root/go/bin

WORKDIR /app

COPY schemas/ ./schemas/
COPY Makefile ./

RUN make gen-perfume-hub

COPY generated/go.mod generated/go.sum ./generated/

COPY services/perfume-hub/go.mod services/perfume-hub/go.sum ./services/perfume-hub/
WORKDIR /app/services/perfume-hub
RUN go mod download

COPY services/perfume-hub/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o perfume ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/perfume-hub/perfume .

EXPOSE 8000 50051

CMD ["./perfume"]
