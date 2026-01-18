PROTO_PATH := schemas/proto
OUT_PATH := shared/go/proto
GO_BIN := $(shell go env GOPATH)/bin
export PATH := $(GO_BIN):$(PATH)

VERSION ?= 1

gen-service:
ifndef SERVICE_NAME
	$(error SERVICE_NAME is required. Usage: make gen-service SERVICE_NAME=perfumist [VERSION=1])
endif
	@rm -rf $(OUT_PATH)/$(SERVICE_NAME)
	@echo "Generating $(SERVICE_NAME) proto (v$(VERSION))..."
	@mkdir -p $(OUT_PATH)
	@protoc \
		--proto_path=$(PROTO_PATH) \
		--go_out=$(OUT_PATH) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_PATH) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/$(SERVICE_NAME)/v$(VERSION)/models/perfume.proto \
		$(PROTO_PATH)/$(SERVICE_NAME)/v$(VERSION)/$(SERVICE_NAME).proto
	@echo "Initializing Go modules..."
	@cd $(OUT_PATH)/$(SERVICE_NAME)/v$(VERSION) && go mod init github.com/zemld/Scently/shared/proto/$(SERVICE_NAME) 2>/dev/null || true
	@cd $(OUT_PATH)/$(SERVICE_NAME)/v$(VERSION) && go mod tidy

gen-perfumist:
	@$(MAKE) gen-service SERVICE_NAME=perfumist VERSION=1

gen-perfume-hub:
	@$(MAKE) gen-service SERVICE_NAME=perfume-hub VERSION=1