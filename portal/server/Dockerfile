# 构建阶段
FROM golang:1.25.0-alpine AS builder

LABEL authors="yeyushilai"
LABEL maintainer="minghe-portal"
LABEL version="1.0.0"

ARG BUILD_ENV=production

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git make gcc musl-dev

# 复制依赖文件
COPY go.mod go.sum ./

# 配置 Go 环境
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0

# 下载依赖
RUN go mod download \
    && go mod verify

# 复制源代码
COPY . .

# 安装代码生成工具
RUN go install github.com/google/wire/cmd/wire@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

# 生成代码
RUN wire ./cmd/server \
    && swag init -g cmd/server/main.go -o docs

# 构建应用
RUN if [ "$BUILD_ENV" = "production" ]; then \
        go build -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=$(git rev-parse HEAD)" -o app ./cmd/server; \
    else \
        go build -o app ./cmd/server; \
    fi

# 运行阶段
FROM alpine:3.18

LABEL authors="yeyushilai"
LABEL maintainer="minghe-portal"

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates wget bash \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config.yaml /app/config.yaml

# 创建日志目录
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8088

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8088/health || exit 1

# 入口点
ENTRYPOINT ["/app/app"]

# 信号处理：确保应用可以优雅关闭
STOPSIGNAL SIGTERM
