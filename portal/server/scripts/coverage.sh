#!/bin/bash

# MingHe Portal 测试覆盖率脚本

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

# 覆盖率阈值
COVERAGE_THRESHOLD=80

# 运行测试并生成覆盖率报告
generate_coverage() {
    print_step "运行测试并生成覆盖率报告..."
    print_info "覆盖率目标: ${COVERAGE_THRESHOLD}%"

    # 运行测试并生成覆盖率文件
    go test -v -coverprofile=coverage.out -covermode=atomic ./...

    # 生成 HTML 覆盖率报告
    go tool cover -html=coverage.out -o coverage.html

    print_info "覆盖率报告生成完成: coverage.html"

    # 获取总覆盖率
    local total_coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

    echo ""
    print_info "总覆盖率: ${total_coverage}%"

    # 检查覆盖率是否达到目标
    local coverage_int=${total_coverage%.*}
    if [ $coverage_int -ge $COVERAGE_THRESHOLD ]; then
        print_info "✓ 覆盖率达到目标 (${COVERAGE_THRESHOLD}%)"
        return 0
    else
        print_warn "✗ 覆盖率未达到目标 (${COVERAGE_THRESHOLD}%)"
        return 1
    fi
}

# 显示每个包的覆盖率
show_package_coverage() {
    print_step "显示各包覆盖率..."
    echo ""
    echo "包覆盖率详情:"
    echo "----------------------------------------"

    go tool cover -func=coverage.out | grep -v "^total:" | while read line; do
        local package=$(echo "$line" | awk '{print $1}' | sed 's/.*\///')
        local coverage=$(echo "$line" | awk '{print $3}')

        # 格式化输出
        if [[ "$coverage" =~ ^[0-9]+\.?[0-9]*%$ ]]; then
            local coverage_num=${coverage%.*}
            if [ $coverage_num -ge 80 ]; then
                echo -e "${GREEN}✓${NC} $package: $coverage"
            elif [ $coverage_num -ge 50 ]; then
                echo -e "${YELLOW}⚠${NC} $package: $coverage"
            else
                echo -e "${RED}✗${NC} $package: $coverage"
            fi
        fi
    done

    echo "----------------------------------------"
}

# 显示未覆盖的代码
show_uncovered() {
    print_step "显示未覆盖的代码..."

    if [ -f "uncovered.out" ]; then
        rm -f "uncovered.out"
    fi

    go tool cover -func=coverage.out | awk '$3 != "100.0%" && $1 != "total:" {print $1, $3}' > uncovered.out

    if [ -s "uncovered.out" ]; then
        print_warn "以下代码未被测试覆盖:"
        cat uncovered.out
    else
        print_info "所有代码都被测试覆盖"
    fi

    rm -f "uncovered.out"
}

# 在浏览器中打开覆盖率报告
open_coverage_report() {
    print_step "打开覆盖率报告..."

    if command -v open &> /dev/null; then
        # macOS
        open coverage.html
    elif command -v xdg-open &> /dev/null; then
        # Linux
        xdg-open coverage.html
    elif command -v start &> /dev/null; then
        # Windows
        start coverage.html
    else
        print_warn "无法自动打开覆盖率报告，请手动打开 coverage.html"
    fi
}

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -o, --open          在浏览器中打开覆盖率报告"
    echo "  -d, --detail        显示详细的覆盖率信息"
    echo "  -u, --uncovered     显示未覆盖的代码"
    echo "  -t, --threshold N   设置覆盖率阈值（默认: 80）"
    echo ""
    echo "示例:"
    echo "  $0                  # 生成覆盖率报告"
    echo "  $0 -o               # 生成并打开覆盖率报告"
    echo "  $0 -d               # 显示详细覆盖率信息"
    echo "  $0 -t 90            # 设置覆盖率为 90%"
}

# 主流程
main() {
    local open_report=false
    local show_detail=false
    local show_uncovered_code=false

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -o|--open)
                open_report=true
                shift
                ;;
            -d|--detail)
                show_detail=true
                shift
                ;;
            -u|--uncovered)
                show_uncovered_code=true
                shift
                ;;
            -t|--threshold)
                COVERAGE_THRESHOLD="$2"
                shift 2
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done

    print_info "MingHe Portal 测试覆盖率分析"
    echo ""

    # 生成覆盖率报告
    generate_coverage
    local exit_code=$?

    # 显示详细信息
    if [ "$show_detail" = true ]; then
        echo ""
        show_package_coverage
    fi

    # 显示未覆盖的代码
    if [ "$show_uncovered_code" = true ]; then
        echo ""
        show_uncovered
    fi

    # 打开报告
    if [ "$open_report" = true ]; then
        echo ""
        open_coverage_report
    fi

    echo ""
    if [ $exit_code -eq 0 ]; then
        print_info "覆盖率分析完成"
    else
        print_warn "覆盖率未达到目标，请继续完善测试"
    fi

    exit $exit_code
}

# 执行主流程
main "$@"
