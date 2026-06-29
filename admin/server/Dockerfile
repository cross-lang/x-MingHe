# ==============================
# 构建阶段：使用 golang:alpine
# ==============================
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 1. 复制模块依赖文件（利用 Docker 缓存）
COPY go.mod go.sum ./

# 2. 下载依赖（比 go mod tidy 更高效，且不修改文件）
RUN go mod download

# 3. 安装 swag 命令行工具
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 4. 复制全部源代码
COPY . .

# 5. 打印 Go 环境（可选，用于调试）
RUN go env

# 6. 生成 Swagger 文档（基于 main.go）
RUN swag init -g main.go -o docs

# 7. 构建二进制文件
RUN go build -o server .

# ==============================
# 运行阶段：使用 alpine:latest
# ==============================
FROM alpine:latest

# 设置时区为 Asia/Shanghai
ENV TZ=Asia/Shanghai

# 安装时区数据并配置
RUN apk add --no-cache tzdata && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

# 创建应用目录
WORKDIR /app

# 从构建阶段复制可执行文件和资源
COPY --from=builder /app/server .
COPY --from=builder /app/resource ./resource/

# 暴露服务端口
EXPOSE 8888

# 启动命令
ENTRYPOINT ["./server", "-c", "config.yaml"]