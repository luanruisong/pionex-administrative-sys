# Go 参数
BINARY_NAME=pionex-administrative-sys
MAIN_FILE=main.go
BUILD_DIR=build

# 版本信息
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 默认目标平台
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build clean run test fmt vet tidy help
.PHONY: build-linux build-linux-arm64 build-darwin build-darwin-arm64 build-windows build-all

# 默认目标
all: build

# 本地构建
build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 跨平台构建
build-linux:
	@echo "Building for linux/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

build-linux-arm64:
	@echo "Building for linux/arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)

build-darwin:
	@echo "Building for darwin/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

build-darwin-arm64:
	@echo "Building for darwin/arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)

# 构建所有平台
build-all: build-linux build-linux-arm64 build-darwin build-darwin-arm64
	@echo "All platforms built successfully!"

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# 运行
run:
	go run $(MAIN_FILE)

# 测试
test:
	go test -v ./...

# 格式化
fmt:
	go fmt ./...

# 检查
vet:
	go vet ./...

# 整理依赖
tidy:
	go mod tidy

# 帮助
help:
	@echo "Usage:"
	@echo "  make build              - 构建当前平台"
	@echo "  make build GOOS=linux GOARCH=amd64  - 指定平台构建"
	@echo "  make build-linux        - 构建 Linux amd64"
	@echo "  make build-linux-arm64  - 构建 Linux arm64"
	@echo "  make build-darwin       - 构建 macOS amd64"
	@echo "  make build-darwin-arm64 - 构建 macOS arm64 (Apple Silicon)"
	@echo "  make build-windows      - 构建 Windows amd64"
	@echo "  make build-all          - 构建所有平台"
	@echo "  make clean              - 清理构建产物"
	@echo "  make run                - 运行应用"
	@echo "  make test               - 运行测试"
	@echo "  make fmt                - 格式化代码"
	@echo "  make vet                - 代码检查"
	@echo "  make tidy               - 整理依赖"