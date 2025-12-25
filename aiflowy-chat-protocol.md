

# IFlowy Chat Protocol Specification v1.1

* **Protocol Name:** `aiflowy-chat-sse`
* **Version:** `1.1`
* **Status:** Draft / Recommended
* **Transport:** Server-Sent Events (SSE)
* **Encoding:** UTF-8



## 1. 设计背景与目标

本协议用于描述 **AIFlowy 对话系统中的服务端事件流通信规范**，支持：

* AI 对话的 **流式输出**
* 模型 **思考过程（Thinking）**
* **工具调用（Tool Calling）**
* **系统 / 业务错误**
* **工作流 / Agent 状态**
* **对话中的用户交互（表单、确认等）**
* **中断与恢复（Suspend / Resume）**

设计目标：

* 前后端解耦
* 协议长期可扩展
* 不绑定具体模型厂商
* 易于与 Workflow / Agent / Chain 架构集成



## 2. 传输层规范（Transport）

* 使用 HTTP + SSE（支持未来扩展为其他协议，比如 WebSocket 等）
* Response Header：

  ```http
  Content-Type: text/event-stream
  Cache-Control: no-cache
  Connection: keep-alive
  ```
* 通信方向：**Server → Client**
* 所有业务数据通过 `data` 字段传输，格式为 **JSON 字符串**



## 3. SSE Event 级别规范

### 3.1 Event Name（固定）

| event   | 含义    |
| - |-------|
| message | 正常业务事件 |
| error   | 错误事件  |
| done    | 流结束事件 |

> ⚠️ **禁止在 event name 中承载业务语义**



## 4. 统一 Envelope 结构（核心）

### 4.1 基本结构

```json
{
  "protocol": "aiflowy-chat-sse",
  "version": "1.1",
  "domain": "llm | tool | system | business | workflow | interaction | debug",
  "type": "string",
  "conversation_id": "string",
  "message_id": "string",
  "index": 0,
  "payload": {},
  "meta": {}
}
```



### 4.2 字段说明

| 字段              | 类型      | 必填 | 说明                     |
| -- |---------| -- |------------------------|
| protocol        | string  | ✔  | 固定值 `aiflowy-chat-sse` |
| version         | string  | ✔  | 协议版本                   |
| domain          | string  | ✔  | 事件所属领域                 |
| type            | string  | ✔  | 领域内事件类型                |
| conversation_id | string  | ✔  | 会话唯一标识                 |
| message_id      | string  | ✖  | assistant 消息 ID        |
| index           | number  | ✖  | 流式输出序号                 |
| payload         | object  | ✔  | 事件数据                   |
| meta            | object  | ✖  | 元信息（token、耗时等）         |



## 5. Domain 定义

| Domain      | 说明             |
| -- | -- |
| llm         | 模型语义输出         |
| tool        | 工具调用与结果        |
| system      | 系统级事件          |
| business    | 业务规则           |
| workflow    | 工作流 / Agent 状态 |
| interaction | 用户交互（表单等）      |
| debug       | 调试信息           |



## 6. llm Domain

### 6.1 thinking

表示模型的思考过程。

#### 流式输出（delta）

```json
{
  "domain": "llm",
  "type": "thinking",
  "payload": {
    "delta": "分析用户需求"
  }
}
```

#### 完整输出（可选）

```json
{
  "domain": "llm",
  "type": "message",
  "payload": {
    "content": "这是一个完整的回答"
  }
}
```


### 6.2 message

#### 流式输出（delta）

```json
{
  "domain": "llm",
  "type": "message",
  "index": 12,
  "payload": {
    "delta": "这是一个"
  }
}
```

#### 完整输出（可选）

```json
{
  "domain": "llm",
  "type": "message",
  "payload": {
    "content": "这是一个完整的回答"
  }
}
```



## 7. tool Domain

### 7.1 tool_call

```json
{
  "domain": "tool",
  "type": "tool_call",
  "payload": {
    "tool_call_id": "call_1",
    "name": "search",
    "arguments": {
      "query": "SSE 协议设计"
    }
  }
}
```



### 7.2 tool_result

```json
{
  "domain": "tool",
  "type": "tool_result",
  "payload": {
    "tool_call_id": "call_1",
    "status": "success | error",
    "result": {}
  }
}
```



## 8. system Domain

### 8.1 error

```json
{
  "domain": "system",
  "type": "error",
  "payload": {
    "code": "MODEL_CONFIG_INVALID",
    "message": "模型配置错误",
    "retryable": false,
    "detail": {}
  }
}
```



### 8.2 status

```json
{
  "domain": "system",
  "type": "status",
  "payload": {
    "state": "initializing | running | suspended | resumed"
  }
}
```



## 9. business Domain

```json
{
  "domain": "business",
  "type": "error",
  "payload": {
    "code": "QUOTA_EXCEEDED",
    "message": "配额不足"
  }
}
```



## 10. workflow Domain

```json
{
  "domain": "workflow",
  "type": "status",
  "payload": {
    "node_id": "node_1",
    "state": "start | suspend | resume | end",
    "reason": "interaction"
  }
}
```



## 11. interaction Domain（对话内交互）

### 11.1 form_request

表示请求用户填写表单，对话进入挂起状态。

```json
{
  "domain": "interaction",
  "type": "form_request",
  "payload": {
    "form_id": "user_info_form",
    "title": "补充信息",
    "description": "请填写以下信息以继续",
    "schema": {
      "type": "object",
      "required": ["age", "email"],
      "properties": {
        "age": {
          "type": "number",
          "title": "年龄"
        },
        "email": {
          "type": "string",
          "title": "邮箱",
          "format": "email"
        }
      }
    },
    "ui": {
      "submit_text": "继续",
      "cancel_text": "取消"
    }
  }
}
```

> 表单 schema **符合 JSON Schema 标准**



### 11.2 form_cancel

```json
{
  "domain": "interaction",
  "type": "form_cancel",
  "payload": {
    "form_id": "user_info_form"
  }
}
```



## 12. 表单提交与恢复（非 SSE）

表单提交通过 **普通 HTTP / WebSocket 请求**：

```json
{
  "conversation_id": "conv_1",
  "form_id": "user_info_form",
  "values": {
    "age": 30,
    "email": "a@b.com"
  }
}
```

成功后服务端恢复 SSE 流。



## 13. done 事件（流结束）

```json
{
  "domain": "system",
  "type": "done",
  "meta": {
    "prompt_tokens": 1234,
    "completion_tokens": 456,
    "latency_ms": 2300
  }
}
```



## 14. 错误处理规则

* 收到 `event: error` 后客户端应终止流
* 错误语义由：

  ```
  domain + type + payload.code
  ```

  共同决定



## 15. 状态机视角（推荐）

```text
RUNNING
  ↓
LLM_OUTPUT
  ↓
INTERACTION_REQUESTED
  ↓
SUSPENDED
  ↓
FORM_SUBMITTED
  ↓
RESUMED
  ↓
RUNNING
  ↓
DONE
```



## 16. 扩展与兼容规则

1. 可新增 domain
2. 可新增 type
3. 不允许删除已有字段
4. payload 可自由扩展
5. 1.x 版本保持向后兼容



## 17. 设计原则

> * SSE 只负责事件流
> * domain 定义责任边界
> * type 定义语义动作
> * payload 定义数据结构
> * 前端不依赖 event name 判断业务，不依赖协议本身，支持其他协议的扩展




