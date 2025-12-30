#!/bin/bash
# AIFlowy Go 后端性能测试脚本
# 使用 wrk 或 ab 进行压测

set -e

BASE_URL="${BASE_URL:-http://localhost:8213}"
API_URL="$BASE_URL/api/v1"
DURATION="${DURATION:-10s}"
THREADS="${THREADS:-4}"
CONNECTIONS="${CONNECTIONS:-100}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo -e "${BLUE}AIFlowy Go 后端性能测试${NC}"
echo "=========================================="
echo "Base URL: $BASE_URL"
echo "Duration: $DURATION"
echo "Threads: $THREADS"
echo "Connections: $CONNECTIONS"
echo ""

# 检查是否安装了压测工具
check_tool() {
    if command -v wrk &> /dev/null; then
        echo "使用 wrk 进行测试"
        USE_WRK=1
    elif command -v ab &> /dev/null; then
        echo "使用 ab (Apache Benchmark) 进行测试"
        USE_WRK=0
    else
        echo -e "${RED}错误: 请安装 wrk 或 ab (Apache Benchmark)${NC}"
        echo "  macOS: brew install wrk"
        echo "  Ubuntu: apt install apache2-utils"
        exit 1
    fi
}

# 获取 Token
get_token() {
    echo -e "\n${YELLOW}[0] 获取认证 Token${NC}"
    TOKEN=$(curl -s "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"account":"admin","password":"123456"}' | \
        python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null || echo "")

    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
        echo -e "${RED}获取 Token 失败${NC}"
        exit 1
    fi
    echo -e "${GREEN}Token 获取成功${NC}"
}

# 运行压测 (wrk)
run_wrk_test() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local body="$4"

    echo -e "\n${YELLOW}[$name]${NC} $method $endpoint"

    if [ "$method" = "GET" ]; then
        wrk -t$THREADS -c$CONNECTIONS -d$DURATION \
            -H "aiflowy-token: $TOKEN" \
            "$BASE_URL$endpoint" 2>&1 | tail -8
    else
        wrk -t$THREADS -c$CONNECTIONS -d$DURATION \
            -s /dev/stdin \
            -H "Content-Type: application/json" \
            -H "aiflowy-token: $TOKEN" \
            "$BASE_URL$endpoint" <<EOF
wrk.method = "POST"
wrk.body   = '$body'
EOF
    fi
}

# 运行压测 (ab)
run_ab_test() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local body="$4"
    local requests=1000

    echo -e "\n${YELLOW}[$name]${NC} $method $endpoint"

    if [ "$method" = "GET" ]; then
        ab -n $requests -c $CONNECTIONS \
            -H "aiflowy-token: $TOKEN" \
            "$BASE_URL$endpoint" 2>&1 | grep -E "(Requests per second|Time per request|Failed requests)"
    else
        echo "$body" > /tmp/ab_body.json
        ab -n $requests -c $CONNECTIONS \
            -T "application/json" \
            -H "aiflowy-token: $TOKEN" \
            -p /tmp/ab_body.json \
            "$BASE_URL$endpoint" 2>&1 | grep -E "(Requests per second|Time per request|Failed requests)"
        rm -f /tmp/ab_body.json
    fi
}

# 选择测试方法
run_test() {
    if [ "$USE_WRK" = "1" ]; then
        run_wrk_test "$@"
    else
        run_ab_test "$@"
    fi
}

check_tool
get_token

echo -e "\n${BLUE}========== 开始压测 ==========${NC}"

# 1. 健康检查 (最简单的 API)
run_test "健康检查" "GET" "/health" ""

# 2. 获取用户信息
run_test "用户信息" "GET" "/api/v1/auth/getUserInfo" ""

# 3. 用户列表 (分页)
run_test "用户列表" "GET" "/api/v1/sysAccount/page?pageNumber=1&pageSize=10" ""

# 4. Bot 列表
run_test "Bot列表" "GET" "/api/v1/bot/list" ""

# 5. 模型列表
run_test "模型列表" "GET" "/api/v1/model/getList" ""

echo -e "\n${BLUE}========== 压测完成 ==========${NC}"
echo ""
echo "性能测试说明:"
echo "- Requests/sec: 每秒请求数 (QPS)"
echo "- Latency: 延迟 (越低越好)"
echo "- Transfer/sec: 每秒传输量"
echo ""
echo "建议基准:"
echo "- 健康检查 API: > 10000 QPS"
echo "- 简单 CRUD API: > 5000 QPS"
echo "- 复杂查询 API: > 1000 QPS"
