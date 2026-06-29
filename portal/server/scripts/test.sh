#!/bin/bash

# MingHe Portal 测试脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 运行测试
run_tests() {
    print_step "运行单元测试..."
    if go test -v -race -count=1 ./...; then
        print_info "测试通过"
        return 0
    else
        print_error "测试失败"
        return 1
    fi
}

# 运行短测试
run_short_tests() {
    print_step "运行短模式测试..."
    if go test -v -short -count=1 ./...; then
        print_info "测试通过"
        return 0
    else
        print_error "测试失败"
        return 1
    fi
}

# 运行特定包的测试
run_package_tests() {
    local package=$1
    print_step "运行包测试: $package"
    if go test -v -count=1 "./$package"; then
        print_info "测试通过"
        return 0
    else
        print_error "测试失败"
        return 1
    fi
}

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help      显示帮助信息"
    echo "  -s, --short     运行短模式测试"
    echo "  -p, --package   运行指定包的测试"
    echo "  -v, --verbose   详细输出"
    echo ""
    echo "示例:"
    echo "  $0              # 运行所有测试"
    echo "  $0 -s           # 运行短模式测试"
    echo "  $0 -p internal/pkg/log  # 运行特定包的测试"
}

# 主流程
main() {
    local short_mode=false
    local specific_package=""
    local verbose=false

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -s|--short)
                short_mode=true
                shift
                ;;
            -p|--package)
                specific_package="$2"
                shift 2
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done

    print_info "MingHe Portal 测试套件"
    echo ""

    # 运行测试
    if [ "$short_mode" = true ]; then
        run_short_tests
    elif [ -n "$specific_package" ]; then
        run_package_tests "$specific_package"
    else
        run_tests
    fi

    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        echo ""
        print_info "所有测试完成"
    else
        echo ""
        print_error "测试失败，请检查错误信息"
    fi

    exit $exit_code
}

# 执行主流程
main "$@"
