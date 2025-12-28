GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m

SCHEMAS_DIR := schemas/proto
GENERATED_DIR := generated/proto
SERVICES_DIR := services

PERFUME_PROTO := $(SCHEMAS_DIR)/models/perfume.proto
REQUESTS_PROTO := $(SCHEMAS_DIR)/models/requests.proto

.PHONY: check-tools
check-tools:
	@echo "$(YELLOW)Checking required tools...$(NC)"
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed"; exit 1; }
	@command -v protoc-gen-go >/dev/null 2>&1 || { echo "Error: protoc-gen-go is not installed. Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "Error: protoc-gen-go-grpc is not installed. Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1; }
	@echo "$(GREEN)✓ All required tools are installed$(NC)"

.PHONY: gen-perfume-hub
gen-perfume-hub: check-tools
	@echo "$(YELLOW)Generating proto for perfume-hub service...$(NC)"
	@mkdir -p $(GENERATED_DIR)/perfume-hub/models
	@mkdir -p $(GENERATED_DIR)/perfume-hub/requests
	@mkdir -p /tmp/proto-gen/perfume-hub/models
	@# Create temporary proto files with correct go_package for perfume-hub
	@sed 's|option go_package = "github.com/zemld/Scently/common/proto/models";|option go_package = "github.com/zemld/Scently/generated/proto/perfume-hub/models";|' $(PERFUME_PROTO) > /tmp/proto-gen/perfume-hub/models/perfume.proto
	@sed 's|option go_package = "github.com/zemld/Scently/common/proto/requests";|option go_package = "github.com/zemld/Scently/generated/proto/perfume-hub/requests";|' $(REQUESTS_PROTO) > /tmp/proto-gen/perfume-hub/models/requests.proto
	@# Generate models/perfume.proto -> models/
	@protoc \
		--proto_path=/tmp/proto-gen/perfume-hub \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfume-hub \
		--go_opt=paths=source_relative \
		models/perfume.proto
	@# Generate models/requests.proto -> requests/
	@protoc \
		--proto_path=/tmp/proto-gen/perfume-hub \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfume-hub/requests \
		--go_opt=paths=source_relative \
		models/requests.proto
	@# Move requests.pb.go from requests/models/ to requests/ if needed
	@if [ -d "$(GENERATED_DIR)/perfume-hub/requests/models" ]; then \
		mv $(GENERATED_DIR)/perfume-hub/requests/models/*.pb.go $(GENERATED_DIR)/perfume-hub/requests/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfume-hub/requests/models 2>/dev/null || true; \
	fi
	@# Move requests.pb.go from requests/requests/ to requests/ if needed
	@if [ -d "$(GENERATED_DIR)/perfume-hub/requests/requests" ]; then \
		mv $(GENERATED_DIR)/perfume-hub/requests/requests/*.pb.go $(GENERATED_DIR)/perfume-hub/requests/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfume-hub/requests/requests 2>/dev/null || true; \
	fi
	@# Remove requests.pb.go from models/ if it was created there
	@rm -f $(GENERATED_DIR)/perfume-hub/models/requests.pb.go
	@# Generate service proto with gRPC (using modified requests.proto)
	@protoc \
		--proto_path=/tmp/proto-gen/perfume-hub \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfume-hub \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GENERATED_DIR)/perfume-hub \
		--go-grpc_opt=paths=source_relative \
		$(SCHEMAS_DIR)/perfume-hub-service.proto
	@# Move generated files to correct location if needed
	@if [ -d "$(GENERATED_DIR)/perfume-hub/models/models" ]; then \
		mv $(GENERATED_DIR)/perfume-hub/models/models/*.pb.go $(GENERATED_DIR)/perfume-hub/models/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfume-hub/models/models 2>/dev/null || true; \
	fi
	@# Cleanup temporary files
	@rm -rf /tmp/proto-gen/perfume-hub
	@echo "$(GREEN)✓ Generated proto for perfume-hub$(NC)"

.PHONY: gen-perfumist
gen-perfumist: check-tools
	@echo "$(YELLOW)Generating proto for perfumist service...$(NC)"
	@mkdir -p $(GENERATED_DIR)/perfumist/models
	@mkdir -p $(GENERATED_DIR)/perfumist/requests
	@mkdir -p /tmp/proto-gen/perfumist/models
	@# Create temporary proto files with correct go_package for perfumist
	@sed 's|option go_package = "github.com/zemld/Scently/common/proto/models";|option go_package = "github.com/zemld/Scently/generated/proto/perfumist/models";|' $(PERFUME_PROTO) > /tmp/proto-gen/perfumist/models/perfume.proto
	@sed 's|option go_package = "github.com/zemld/Scently/common/proto/requests";|option go_package = "github.com/zemld/Scently/generated/proto/perfumist/requests";|' $(REQUESTS_PROTO) > /tmp/proto-gen/perfumist/models/requests.proto
	@# Generate models/perfume.proto -> models/
	@protoc \
		--proto_path=/tmp/proto-gen/perfumist \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfumist \
		--go_opt=paths=source_relative \
		models/perfume.proto
	@# Generate models/requests.proto -> requests/
	@protoc \
		--proto_path=/tmp/proto-gen/perfumist \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfumist/requests \
		--go_opt=paths=source_relative \
		models/requests.proto
	@# Move requests.pb.go from requests/models/ to requests/ if needed
	@if [ -d "$(GENERATED_DIR)/perfumist/requests/models" ]; then \
		mv $(GENERATED_DIR)/perfumist/requests/models/*.pb.go $(GENERATED_DIR)/perfumist/requests/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfumist/requests/models 2>/dev/null || true; \
	fi
	@# Move requests.pb.go from requests/requests/ to requests/ if needed
	@if [ -d "$(GENERATED_DIR)/perfumist/requests/requests" ]; then \
		mv $(GENERATED_DIR)/perfumist/requests/requests/*.pb.go $(GENERATED_DIR)/perfumist/requests/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfumist/requests/requests 2>/dev/null || true; \
	fi
	@# Remove requests.pb.go from models/ if it was created there
	@rm -f $(GENERATED_DIR)/perfumist/models/requests.pb.go
	@# Generate service proto with gRPC (using modified requests.proto)
	@protoc \
		--proto_path=/tmp/proto-gen/perfumist \
		--proto_path=$(SCHEMAS_DIR) \
		--go_out=$(GENERATED_DIR)/perfumist \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GENERATED_DIR)/perfumist \
		--go-grpc_opt=paths=source_relative \
		$(SCHEMAS_DIR)/perfumist-service.proto
	@# Move generated files to correct location if needed
	@if [ -d "$(GENERATED_DIR)/perfumist/models/models" ]; then \
		mv $(GENERATED_DIR)/perfumist/models/models/*.pb.go $(GENERATED_DIR)/perfumist/models/ 2>/dev/null || true; \
		rmdir $(GENERATED_DIR)/perfumist/models/models 2>/dev/null || true; \
	fi
	@# Cleanup temporary files
	@rm -rf /tmp/proto-gen/perfumist
	@echo "$(GREEN)✓ Generated proto for perfumist$(NC)"

.PHONY: gen-ai-advisor
gen-ai-advisor:
	@echo "$(YELLOW)Generating proto for ai-advisor service...$(NC)"
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed"; exit 1; }
	@mkdir -p $(GENERATED_DIR)/ai-advisor
	@protoc \
		--proto_path=$(SCHEMAS_DIR) \
		--python_out=$(GENERATED_DIR)/ai-advisor \
		$(SCHEMAS_DIR)/models/perfume.proto \
		$(SCHEMAS_DIR)/models/requests.proto \
		$(SCHEMAS_DIR)/ai-service.proto
	@if command -v grpc_python_plugin >/dev/null 2>&1 || command -v protoc-gen-grpc_python >/dev/null 2>&1; then \
		protoc \
			--proto_path=$(SCHEMAS_DIR) \
			--grpc_python_out=$(GENERATED_DIR)/ai-advisor \
			$(SCHEMAS_DIR)/ai-service.proto; \
		echo "$(GREEN)✓ Generated gRPC code for ai-advisor$(NC)"; \
	else \
		echo "$(YELLOW)⚠ grpc_python_plugin not found, skipping gRPC code generation$(NC)"; \
		echo "$(YELLOW)  Install with: pip install grpcio-tools$(NC)"; \
	fi
	@echo "$(GREEN)✓ Generated proto for ai-advisor$(NC)"

.PHONY: gen-gateway
gen-gateway: check-tools
	@echo "$(YELLOW)Generating proto for gateway service...$(NC)"
	@mkdir -p $(GENERATED_DIR)/gateway
	@echo "$(GREEN)✓ Gateway proto generation (if needed)$(NC)"

.PHONY: gen-all
gen-all: gen-perfume-hub gen-perfumist gen-ai-advisor
	@echo "$(GREEN)✓ All proto files generated$(NC)"

.PHONY: clean-proto
clean-proto:
	@echo "$(YELLOW)Cleaning generated proto files...$(NC)"
	@rm -rf $(GENERATED_DIR)
	@rm -rf /tmp/proto-gen
	@echo "$(GREEN)✓ Cleaned all generated proto files$(NC)"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make gen-perfume-hub   - Generate proto for perfume-hub service"
	@echo "  make gen-perfumist     - Generate proto for perfumist service"
	@echo "  make gen-ai-advisor    - Generate proto for ai-advisor service"
	@echo "  make gen-gateway       - Generate proto for gateway service"
	@echo "  make gen-all           - Generate proto for all services"
	@echo "  make clean-proto       - Clean all generated proto files"
	@echo "  make check-tools       - Check if required tools are installed"
