#!/bin/bash

# MingHe Portal 开发启动脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Go 是否安装
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 1.25.0 或更高版本"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go 版本: $GO_VERSION"
}

# 检查 Air 是否安装
check_air() {
    if ! command -v air &> /dev/null; then
        print_warn "Air 未安装，正在安装..."
        go install github.com/cosmtrek/air@latest
    fi
    print_info "Air 已就绪"
}

# 检查 Wire 是否安装
check_wire() {
    if ! command -v wire &> /dev/null; then
        print_warn "Wire 未安装，正在安装..."
        go install github.com/google/wire/cmd/wire@latest
    fi
    print_info "Wire 已就绪"
}

# 检查配置文件
check_config() {
    if [ ! -f "config.yaml" ] && [ ! -f "config_local.yaml" ]; then
        print_warn "未找到配置文件，将使用 config.yaml"
    fi

    # 如果有 config_local.yaml，使用它
    if [ -f "config_local.yaml" ]; then
        print_info "使用配置文件: config_local.yaml"
        CONFIG_FILE="config_local.yaml"
    else
        print_info "使用配置文件: config.yaml"
        CONFIG_FILE="config.yaml"
    fi

    # 设置环境变量
    export CONFIG_FILE
}

# 生成依赖注入代码
generate_wire() {
    print_info "生成依赖注入代码..."
    if ! wire ./cmd/server; then
        print_error "Wire 代码生成失败"
        exit 1
    fi
    print_info "Wire 代码生成完成"
}

# 启动开发服务器
start_dev() {
    print_info "启动开发服务器..."
    air
}

# 主流程
main() {
    print_info "MingHe Portal 开发环境启动"
    echo ""

    # 检查环境
    check_go
    check_air
    check_wire

    # 检查配置
    check_config

    # 生成代码
    generate_wire

    # 启动服务
    echo ""
    print_info "准备启动开发服务器..."
    echo ""
    start_dev
}

# 捕获中断信号
trap 'print_info "开发服务器已停止"; exit 0' INT TERM

# 执行主流程
main
