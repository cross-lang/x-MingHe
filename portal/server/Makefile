.PHONY: help dev test build docker clean install-tools wire swagger run

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make dev          - 启动热重载开发服务器"
	@echo "  make run          - 运行应用（不热重载）"
	@echo "  make build        - 构建应用"
	@echo "  make test         - 运行所有测试"
	@echo "  make coverage     - 生成测试覆盖率报告"
	@echo "  make wire         - 生成依赖注入代码"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make docker       - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"
	@echo "  make clean        - 清理构建文件"
	@echo "  make install-tools- 安装开发工具"

# 安装开发工具
install-tools:
	@echo "安装开发工具..."
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/cosmtrek/air@latest
	@echo "开发工具安装完成"

# 生成依赖注入代码
wire:
	@echo "生成依赖注入代码..."
	@wire ./cmd/server
	@echo "依赖注入代码生成完成"

# 生成 Swagger 文档
swagger:
	@echo "生成 Swagger 文档..."
	@swag init -g cmd/server/main.go -o docs
	@echo "Swagger 文档生成完成"

# 生成代码（wire + swagger）
gen: wire swagger

# 启动热重载开发服务器
dev:
	@echo "启动热重载开发服务器..."
	@air

# 运行应用
run:
	@echo "运行应用..."
	@go run ./cmd/server

# 构建应用
build:
	@echo "构建应用..."
	@go build -o ./bin/portal ./cmd/server
	@echo "构建完成: ./bin/portal"

# 构建应用（Windows）
build-windows:
	@echo "构建 Windows 应用..."
	@SET GOOS=windows
	@SET GOARCH=amd64
	@go build -o ./bin/portal.exe ./cmd/server
	@echo "Windows 构建完成: ./bin/portal.exe"

# 构建应用（Linux）
build-linux:
	@echo "构建 Linux 应用..."
	@SET GOOS=linux
	@SET GOARCH=amd64
	@go build -o ./bin/portal-linux ./cmd/server
	@echo "Linux 构建完成: ./bin/portal-linux"

# 构建应用（macOS）
build-darwin:
	@echo "构建 macOS 应用..."
	@SET GOOS=darwin
	@SET GOARCH=amd64
	@go build -o ./bin/portal-darwin ./cmd/server
	@echo "macOS 构建完成: ./bin/portal-darwin"

# 构建所有平台
build-all: build-windows build-linux build-darwin

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告生成完成: coverage.html"

# 运行测试（短模式）
test-short:
	@echo "运行短模式测试..."
	@go test -short -v ./...

# 运行基准测试
benchmark:
	@echo "运行基准测试..."
	@go test -bench=. -benchmem ./...

# 检查代码
lint:
	@echo "检查代码..."
	@golangci-lint run ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 检查依赖
mod-tidy:
	@echo "整理依赖..."
	@go mod tidy
	@go mod verify

# 更新依赖
mod-update:
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy

# 构建 Docker 镜像
docker:
	@echo "构建 Docker 镜像..."
	@docker build -t minghe-portal:latest .
	@echo "Docker 镜像构建完成: minghe-portal:latest"

# 运行 Docker 容器
docker-run:
	@echo "运行 Docker 容器..."
	@docker run -p 8088:8088 --name minghe-portal minghe-portal:latest

# 停止 Docker 容器
docker-stop:
	@echo "停止 Docker 容器..."
	@docker stop minghe-portal || true
	@docker rm minghe-portal || true

# 重启 Docker 容器
docker-restart: docker-stop docker-run

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf ./bin
	@rm -rf ./tmp
	@rm -f coverage.out coverage.html build-errors.log
	@echo "清理完成"

# 清理所有（包括 Docker）
clean-all: clean docker-stop
	@echo "清理所有文件..."
	@docker rmi minghe-portal:latest || true
	@echo "全部清理完成"

# 完整构建流程（wire -> swagger -> build）
full-build: gen build
	@echo "完整构建完成"

# 完整测试流程（test -> coverage）
full-test: test coverage
	@echo "完整测试完成"

# 初始化项目
init: install-tools mod-tidy gen
	@echo "项目初始化完成"
