.PHONY: run build clean test mod-tidy mod-download docker-build docker-run help

# 项目名称
PROJECT_NAME=go-demo
# 输出的二进制文件名
BINARY_NAME=app

# 默认目标
.DEFAULT_GOAL := help

# 运行项目
run:
	@echo "🚀 运行项目..."
	go run main.go

# 编译项目
build:
	@echo "🔨 编译项目..."
	go build -o ${BINARY_NAME} main.go
	@echo "✅ 编译完成: ${BINARY_NAME}"

# 编译 Linux 版本
build-linux:
	@echo "🔨 编译 Linux 版本..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux main.go
	@echo "✅ 编译完成: ${BINARY_NAME}-linux"

# 清理编译文件
clean:
	@echo "🧹 清理编译文件..."
	@if exist ${BINARY_NAME}.exe del ${BINARY_NAME}.exe
	@if exist ${BINARY_NAME} del ${BINARY_NAME}
	@if exist ${BINARY_NAME}-linux del ${BINARY_NAME}-linux
	@echo "✅ 清理完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 整理依赖
mod-tidy:
	@echo "📦 整理依赖..."
	go mod tidy
	@echo "✅ 依赖整理完成"

# 下载依赖
mod-download:
	@echo "📥 下载依赖..."
	go mod download
	@echo "✅ 依赖下载完成"

# 查看依赖
mod-verify:
	@echo "🔍 验证依赖..."
	go mod verify

# 热重载（需要安装 air: go install github.com/air-verse/air@latest）
dev:
	@echo "🔥 启动热重载..."
	air

# 代码格式化
fmt:
	@echo "📝 格式化代码..."
	go fmt ./...
	@echo "✅ 格式化完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	golangci-lint run
	@echo "✅ 检查完成"

# 生成 Swagger 文档
swagger:
	@echo "📝 生成 Swagger 文档..."
	swag init
	@echo "✅ Swagger 文档生成完成"
	@echo "访问: http://localhost:8080/swagger/index.html"

# Docker 构建
docker-build:
	@echo "🐳 构建 Docker 镜像..."
	docker build -t ${PROJECT_NAME}:latest .
	@echo "✅ Docker 镜像构建完成"

# Docker 运行
docker-run:
	@echo "🐳 运行 Docker 容器..."
	docker run -p 8080:8080 --name ${PROJECT_NAME} ${PROJECT_NAME}:latest

# Docker Compose 启动
docker-up:
	@echo "🐳 启动 Docker Compose..."
	docker-compose up -d
	@echo "✅ 服务启动成功"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

# Docker Compose 停止
docker-down:
	@echo "🐳 停止 Docker Compose..."
	docker-compose down
	@echo "✅ 服务已停止"

# Docker Compose 重启
docker-restart:
	@echo "🐳 重启 Docker Compose..."
	docker-compose restart
	@echo "✅ 服务已重启"

# Docker Compose 日志
docker-logs:
	@echo "📋 查看 Docker Compose 日志..."
	docker-compose logs -f

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make run           - 运行项目"
	@echo "  make build         - 编译项目"
	@echo "  make build-linux   - 编译 Linux 版本"
	@echo "  make clean         - 清理编译文件"
	@echo "  make test          - 运行测试"
	@echo "  make mod-tidy      - 整理依赖"
	@echo "  make mod-download  - 下载依赖"
	@echo "  make mod-verify    - 验证依赖"
	@echo "  make dev           - 热重载开发 (需要安装 air)"
	@echo "  make fmt           - 格式化代码"
	@echo "  make lint          - 代码检查 (需要安装 golangci-lint)"
	@echo "  make swagger       - 生成 Swagger 文档 (需要安装 swag)"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-run    - 运行 Docker 容器"
	@echo "  make docker-up     - 启动 Docker Compose"
	@echo "  make docker-down   - 停止 Docker Compose"
	@echo "  make docker-restart- 重启 Docker Compose"
	@echo "  make docker-logs   - 查看 Docker Compose 日志"
	@echo "  make help          - 显示帮助信息"

