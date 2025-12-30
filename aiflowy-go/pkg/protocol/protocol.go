package protocol

import "encoding/json"

const (
	// Protocol info
	ProtocolName    = "aiflowy-chat"
	ProtocolVersion = "1.1"
)

// Domain constants
const (
	DomainLLM         = "llm"
	DomainTool        = "tool"
	DomainSystem      = "system"
	DomainBusiness    = "business"
	DomainWorkflow    = "workflow"
	DomainInteraction = "interaction"
	DomainDebug       = "debug"
)

// LLM domain types
const (
	TypeMessage  = "message"
	TypeThinking = "thinking"
)

// Tool domain types
const (
	TypeToolCall   = "tool_call"
	TypeToolResult = "tool_result"
)

// System domain types
const (
	TypeError  = "error"
	TypeStatus = "status"
	TypeDone   = "done"
)

// Business domain types
const (
	TypeBusinessError = "error"
)

// Workflow domain types
const (
	TypeWorkflowStatus = "status"
)

// Interaction domain types
const (
	TypeFormRequest = "form_request"
	TypeFormCancel  = "form_cancel"
)

// SSE Event names
const (
	EventMessage = "message"
	EventError   = "error"
	EventDone    = "done"
)

// Envelope is the unified message envelope for AIFlowy Chat Protocol v1.1
type Envelope struct {
	Protocol       string      `json:"protocol"`
	Version        string      `json:"version"`
	Domain         string      `json:"domain"`
	Type           string      `json:"type"`
	ConversationID string      `json:"conversation_id,omitempty"`
	MessageID      string      `json:"message_id,omitempty"`
	Index          int         `json:"index,omitempty"`
	Payload        interface{} `json:"payload"`
	Meta           *Meta       `json:"meta,omitempty"`
}

// Meta contains metadata like token usage and timing
type Meta struct {
	PromptTokens     int    `json:"prompt_tokens,omitempty"`
	CompletionTokens int    `json:"completion_tokens,omitempty"`
	TotalTokens      int    `json:"total_tokens,omitempty"`
	LatencyMs        int64  `json:"latency_ms,omitempty"`
	ModelName        string `json:"model_name,omitempty"`
	FinishReason     string `json:"finish_reason,omitempty"`
}

// LLM Payloads

// MessagePayload for llm.message and llm.thinking
type MessagePayload struct {
	Delta   string `json:"delta,omitempty"`   // For streaming delta
	Content string `json:"content,omitempty"` // For full content
}

// Tool Payloads

// ToolCallPayload for tool.tool_call
type ToolCallPayload struct {
	ToolCallID string                 `json:"tool_call_id"`
	Name       string                 `json:"name"`
	Arguments  map[string]interface{} `json:"arguments"`
}

// ToolResultPayload for tool.tool_result
type ToolResultPayload struct {
	ToolCallID string      `json:"tool_call_id"`
	Status     string      `json:"status"` // success or error
	Result     interface{} `json:"result"`
}

// System Payloads

// ErrorPayload for system.error and business.error
type ErrorPayload struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Retryable bool        `json:"retryable,omitempty"`
	Detail    interface{} `json:"detail,omitempty"`
}

// StatusPayload for system.status
type StatusPayload struct {
	State string `json:"state"` // initializing, running, suspended, resumed
}

// Workflow Payloads

// WorkflowStatusPayload for workflow.status
type WorkflowStatusPayload struct {
	NodeID string `json:"node_id,omitempty"`
	State  string `json:"state"`  // start, suspend, resume, end
	Reason string `json:"reason,omitempty"`
}

// Builder provides convenient methods to build envelopes
type Builder struct {
	conversationID string
	messageID      string
	index          int
}

// NewBuilder creates a new envelope builder
func NewBuilder(conversationID, messageID string) *Builder {
	return &Builder{
		conversationID: conversationID,
		messageID:      messageID,
		index:          0,
	}
}

// NextIndex increments and returns the next index
func (b *Builder) NextIndex() int {
	b.index++
	return b.index
}

// Reset resets the index counter
func (b *Builder) Reset() {
	b.index = 0
}

// newEnvelope creates a base envelope
func (b *Builder) newEnvelope(domain, typ string, payload interface{}) *Envelope {
	return &Envelope{
		Protocol:       ProtocolName,
		Version:        ProtocolVersion,
		Domain:         domain,
		Type:           typ,
		ConversationID: b.conversationID,
		MessageID:      b.messageID,
		Payload:        payload,
	}
}

// LLMMessageDelta creates a llm.message envelope with delta content
func (b *Builder) LLMMessageDelta(delta string) *Envelope {
	env := b.newEnvelope(DomainLLM, TypeMessage, &MessagePayload{Delta: delta})
	env.Index = b.NextIndex()
	return env
}

// LLMThinkingDelta creates a llm.thinking envelope with delta content
func (b *Builder) LLMThinkingDelta(delta string) *Envelope {
	env := b.newEnvelope(DomainLLM, TypeThinking, &MessagePayload{Delta: delta})
	env.Index = b.NextIndex()
	return env
}

// LLMMessageFull creates a llm.message envelope with full content
func (b *Builder) LLMMessageFull(content string) *Envelope {
	return b.newEnvelope(DomainLLM, TypeMessage, &MessagePayload{Content: content})
}

// ToolCall creates a tool.tool_call envelope
func (b *Builder) ToolCall(callID, name string, args map[string]interface{}) *Envelope {
	return b.newEnvelope(DomainTool, TypeToolCall, &ToolCallPayload{
		ToolCallID: callID,
		Name:       name,
		Arguments:  args,
	})
}

// ToolResult creates a tool.tool_result envelope
func (b *Builder) ToolResult(callID, status string, result interface{}) *Envelope {
	return b.newEnvelope(DomainTool, TypeToolResult, &ToolResultPayload{
		ToolCallID: callID,
		Status:     status,
		Result:     result,
	})
}

// SystemError creates a system.error envelope
func (b *Builder) SystemError(code, message string, retryable bool) *Envelope {
	return b.newEnvelope(DomainSystem, TypeError, &ErrorPayload{
		Code:      code,
		Message:   message,
		Retryable: retryable,
	})
}

// SystemStatus creates a system.status envelope
func (b *Builder) SystemStatus(state string) *Envelope {
	return b.newEnvelope(DomainSystem, TypeStatus, &StatusPayload{State: state})
}

// SystemDone creates a system.done envelope
func (b *Builder) SystemDone(meta *Meta) *Envelope {
	env := b.newEnvelope(DomainSystem, TypeDone, nil)
	env.Meta = meta
	return env
}

// BusinessError creates a business.error envelope
func (b *Builder) BusinessError(code, message string) *Envelope {
	return b.newEnvelope(DomainBusiness, TypeBusinessError, &ErrorPayload{
		Code:    code,
		Message: message,
	})
}

// WorkflowStatus creates a workflow.status envelope
func (b *Builder) WorkflowStatus(nodeID, state, reason string) *Envelope {
	return b.newEnvelope(DomainWorkflow, TypeWorkflowStatus, &WorkflowStatusPayload{
		NodeID: nodeID,
		State:  state,
		Reason: reason,
	})
}

// ToJSON converts an envelope to JSON string
func (e *Envelope) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToSSE formats an envelope as SSE data
func (e *Envelope) ToSSE() (string, error) {
	jsonStr, err := e.ToJSON()
	if err != nil {
		return "", err
	}
	return "event: " + EventMessage + "\ndata: " + jsonStr + "\n\n", nil
}
