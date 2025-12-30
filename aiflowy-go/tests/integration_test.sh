#!/bin/bash
# AIFlowy Go 后端集成测试脚本
# 测试所有主要 API 端点

set -e

BASE_URL="${BASE_URL:-http://localhost:8213}"
API_URL="$BASE_URL/api/v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 计数器
PASSED=0
FAILED=0
TOTAL=0

# 测试函数
test_api() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_code="${5:-0}"

    TOTAL=$((TOTAL + 1))

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_URL$endpoint" -H "aiflowy-token: $TOKEN" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "aiflowy-token: $TOKEN" \
            -d "$data" 2>/dev/null)
    fi

    http_code=$(echo "$response" | tail -1)
    body=$(echo "$response" | sed '$d')
    code=$(echo "$body" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code', -1))" 2>/dev/null || echo "-1")

    if [ "$code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $name"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $name (expected code=$expected_code, got code=$code, http=$http_code)"
        echo "  Response: $(echo "$body" | head -c 200)"
        FAILED=$((FAILED + 1))
    fi
}

echo "=========================================="
echo "AIFlowy Go 后端集成测试"
echo "=========================================="
echo "Base URL: $BASE_URL"
echo ""

# 1. 健康检查
echo -e "${YELLOW}[1] 健康检查${NC}"
TOTAL=$((TOTAL + 1))
health_resp=$(curl -s "$BASE_URL/health" 2>/dev/null)
health_status=$(echo "$health_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('status',''))" 2>/dev/null || echo "")
if [ "$health_status" = "ok" ]; then
    echo -e "${GREEN}✓ PASS${NC}: Health check"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAIL${NC}: Health check"
    FAILED=$((FAILED + 1))
fi

# 2. 登录认证
echo -e "\n${YELLOW}[2] 认证登录${NC}"
TOTAL=$((TOTAL + 1))
login_resp=$(curl -s "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"account":"admin","password":"123456"}' 2>/dev/null)
TOKEN=$(echo "$login_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null || echo "")

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASS${NC}: Login successful"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAIL${NC}: Login failed"
    echo "Response: $login_resp"
    FAILED=$((FAILED + 1))
    echo -e "\n${RED}无法继续测试：登录失败${NC}"
    exit 1
fi

# 3. 认证 API
echo -e "\n${YELLOW}[3] 认证 API${NC}"
test_api "获取用户信息" "GET" "/auth/getUserInfo"
test_api "获取权限列表" "GET" "/auth/getPermissions"
test_api "获取角色列表" "GET" "/auth/getRoles"

# 4. 系统管理 API
echo -e "\n${YELLOW}[4] 系统管理 API${NC}"
test_api "用户列表" "GET" "/sysAccount/page?pageNumber=1&pageSize=10"
test_api "角色列表" "GET" "/sysRole/list"
test_api "菜单列表" "GET" "/sysMenu/list"
test_api "部门列表" "GET" "/sysDept/list"
test_api "字典列表" "GET" "/sysDict/list"

# 5. 模型管理 API
echo -e "\n${YELLOW}[5] 模型管理 API${NC}"
test_api "模型提供商列表" "GET" "/modelProvider/list"
test_api "模型列表" "GET" "/model/getList"

# 6. Bot API
echo -e "\n${YELLOW}[6] Bot API${NC}"
test_api "Bot 列表" "GET" "/bot/list"
test_api "Bot 分类列表" "GET" "/botCategory/list"
test_api "生成会话 ID" "GET" "/bot/generateConversationId"

# 获取第一个 Bot ID
BOT_ID=$(curl -s "$API_URL/bot/list" -H "aiflowy-token: $TOKEN" 2>/dev/null | \
    python3 -c "import sys,json; d=json.load(sys.stdin).get('data',[]); print(d[0]['id'] if d else '')" 2>/dev/null || echo "")

if [ -n "$BOT_ID" ] && [ "$BOT_ID" != "null" ]; then
    test_api "Bot 详情" "GET" "/bot/getDetail?id=$BOT_ID"

    # 非流式聊天测试
    TOTAL=$((TOTAL + 1))
    echo -n "非流式聊天: "
    chat_resp=$(curl -s "$API_URL/bot/chat" \
        -H "Content-Type: application/json" \
        -H "aiflowy-token: $TOKEN" \
        -d "{\"botId\":\"$BOT_ID\",\"message\":\"hello\",\"stream\":false}" 2>/dev/null)
    chat_code=$(echo "$chat_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code', -1))" 2>/dev/null || echo "-1")
    if [ "$chat_code" = "0" ]; then
        echo -e "${GREEN}✓ PASS${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC} (code=$chat_code)"
        FAILED=$((FAILED + 1))
    fi
else
    echo -e "${YELLOW}⚠ SKIP${NC}: 没有可用的 Bot"
fi

# 7. 知识库 API
echo -e "\n${YELLOW}[7] 知识库 API${NC}"
test_api "知识库列表" "GET" "/documentCollection/list"

# 8. 工作流 API
echo -e "\n${YELLOW}[8] 工作流 API${NC}"
test_api "工作流列表" "GET" "/workflow/list"
test_api "工作流分类列表" "GET" "/workflowCategory/list"

# 9. 插件 API
echo -e "\n${YELLOW}[9] 插件 API${NC}"
test_api "插件列表" "POST" "/plugin/getList" "{}"

# 10. 辅助功能 API
echo -e "\n${YELLOW}[10] 辅助功能 API${NC}"
test_api "系统 API 密钥列表" "GET" "/sysApiKey/page?pageNumber=1&pageSize=10"
test_api "操作日志列表" "GET" "/sysLog/page?pageNumber=1&pageSize=10"
test_api "定时任务列表" "GET" "/sysJob/page?pageNumber=1&pageSize=10"
test_api "获取 Cron 下次执行时间" "GET" "/sysJob/getNextTimes?cronExpression=0+*+*+*+*+?"

# 11. 公共 API
echo -e "\n${YELLOW}[11] 公共 API${NC}"
TOTAL=$((TOTAL + 1))
captcha_resp=$(curl -s "$API_URL/public/getCaptcha" 2>/dev/null)
captcha_code=$(echo "$captcha_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code', -1))" 2>/dev/null || echo "-1")
if [ "$captcha_code" = "0" ]; then
    echo -e "${GREEN}✓ PASS${NC}: 获取验证码"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAIL${NC}: 获取验证码 (code=$captcha_code)"
    FAILED=$((FAILED + 1))
fi

# 汇总结果
echo ""
echo "=========================================="
echo "测试结果汇总"
echo "=========================================="
echo -e "总计: $TOTAL"
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"

if [ $FAILED -gt 0 ]; then
    echo -e "\n${RED}测试未通过${NC}"
    exit 1
else
    echo -e "\n${GREEN}所有测试通过！${NC}"
    exit 0
fi
