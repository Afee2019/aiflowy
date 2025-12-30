package protocol

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	// Test protocol info
	if ProtocolName != "aiflowy-chat" {
		t.Errorf("expected protocol name 'aiflowy-chat', got '%s'", ProtocolName)
	}
	if ProtocolVersion != "1.1" {
		t.Errorf("expected protocol version '1.1', got '%s'", ProtocolVersion)
	}

	// Test domain constants
	domains := []string{DomainLLM, DomainTool, DomainSystem, DomainBusiness, DomainWorkflow, DomainInteraction, DomainDebug}
	expectedDomains := []string{"llm", "tool", "system", "business", "workflow", "interaction", "debug"}
	for i, d := range domains {
		if d != expectedDomains[i] {
			t.Errorf("expected domain '%s', got '%s'", expectedDomains[i], d)
		}
	}
}

func TestNewBuilder(t *testing.T) {
	conversationID := "conv-123"
	messageID := "msg-456"

	b := NewBuilder(conversationID, messageID)

	if b.conversationID != conversationID {
		t.Errorf("expected conversationID '%s', got '%s'", conversationID, b.conversationID)
	}
	if b.messageID != messageID {
		t.Errorf("expected messageID '%s', got '%s'", messageID, b.messageID)
	}
	if b.index != 0 {
		t.Errorf("expected index 0, got %d", b.index)
	}
}

func TestBuilderNextIndex(t *testing.T) {
	b := NewBuilder("conv", "msg")

	for i := 1; i <= 5; i++ {
		idx := b.NextIndex()
		if idx != i {
			t.Errorf("expected index %d, got %d", i, idx)
		}
	}
}

func TestBuilderReset(t *testing.T) {
	b := NewBuilder("conv", "msg")
	b.NextIndex()
	b.NextIndex()
	b.Reset()

	if b.index != 0 {
		t.Errorf("expected index 0 after reset, got %d", b.index)
	}
}

func TestLLMMessageDelta(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	delta := "Hello"

	env := b.LLMMessageDelta(delta)

	if env.Protocol != ProtocolName {
		t.Errorf("expected protocol '%s', got '%s'", ProtocolName, env.Protocol)
	}
	if env.Version != ProtocolVersion {
		t.Errorf("expected version '%s', got '%s'", ProtocolVersion, env.Version)
	}
	if env.Domain != DomainLLM {
		t.Errorf("expected domain '%s', got '%s'", DomainLLM, env.Domain)
	}
	if env.Type != TypeMessage {
		t.Errorf("expected type '%s', got '%s'", TypeMessage, env.Type)
	}
	if env.ConversationID != "conv-1" {
		t.Errorf("expected conversationID 'conv-1', got '%s'", env.ConversationID)
	}
	if env.MessageID != "msg-1" {
		t.Errorf("expected messageID 'msg-1', got '%s'", env.MessageID)
	}
	if env.Index != 1 {
		t.Errorf("expected index 1, got %d", env.Index)
	}

	payload, ok := env.Payload.(*MessagePayload)
	if !ok {
		t.Fatalf("expected MessagePayload type")
	}
	if payload.Delta != delta {
		t.Errorf("expected delta '%s', got '%s'", delta, payload.Delta)
	}
}

func TestLLMThinkingDelta(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	delta := "Thinking..."

	env := b.LLMThinkingDelta(delta)

	if env.Domain != DomainLLM {
		t.Errorf("expected domain '%s', got '%s'", DomainLLM, env.Domain)
	}
	if env.Type != TypeThinking {
		t.Errorf("expected type '%s', got '%s'", TypeThinking, env.Type)
	}

	payload, ok := env.Payload.(*MessagePayload)
	if !ok {
		t.Fatalf("expected MessagePayload type")
	}
	if payload.Delta != delta {
		t.Errorf("expected delta '%s', got '%s'", delta, payload.Delta)
	}
}

func TestLLMMessageFull(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	content := "Full message content"

	env := b.LLMMessageFull(content)

	if env.Domain != DomainLLM {
		t.Errorf("expected domain '%s', got '%s'", DomainLLM, env.Domain)
	}
	if env.Type != TypeMessage {
		t.Errorf("expected type '%s', got '%s'", TypeMessage, env.Type)
	}
	if env.Index != 0 {
		t.Errorf("expected index 0 for full message, got %d", env.Index)
	}

	payload, ok := env.Payload.(*MessagePayload)
	if !ok {
		t.Fatalf("expected MessagePayload type")
	}
	if payload.Content != content {
		t.Errorf("expected content '%s', got '%s'", content, payload.Content)
	}
}

func TestToolCall(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	callID := "call-123"
	name := "get_weather"
	args := map[string]interface{}{"city": "Beijing"}

	env := b.ToolCall(callID, name, args)

	if env.Domain != DomainTool {
		t.Errorf("expected domain '%s', got '%s'", DomainTool, env.Domain)
	}
	if env.Type != TypeToolCall {
		t.Errorf("expected type '%s', got '%s'", TypeToolCall, env.Type)
	}

	payload, ok := env.Payload.(*ToolCallPayload)
	if !ok {
		t.Fatalf("expected ToolCallPayload type")
	}
	if payload.ToolCallID != callID {
		t.Errorf("expected callID '%s', got '%s'", callID, payload.ToolCallID)
	}
	if payload.Name != name {
		t.Errorf("expected name '%s', got '%s'", name, payload.Name)
	}
	if payload.Arguments["city"] != "Beijing" {
		t.Errorf("expected city 'Beijing', got '%v'", payload.Arguments["city"])
	}
}

func TestToolResult(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	callID := "call-123"
	status := "success"
	result := map[string]string{"weather": "sunny"}

	env := b.ToolResult(callID, status, result)

	if env.Domain != DomainTool {
		t.Errorf("expected domain '%s', got '%s'", DomainTool, env.Domain)
	}
	if env.Type != TypeToolResult {
		t.Errorf("expected type '%s', got '%s'", TypeToolResult, env.Type)
	}

	payload, ok := env.Payload.(*ToolResultPayload)
	if !ok {
		t.Fatalf("expected ToolResultPayload type")
	}
	if payload.ToolCallID != callID {
		t.Errorf("expected callID '%s', got '%s'", callID, payload.ToolCallID)
	}
	if payload.Status != status {
		t.Errorf("expected status '%s', got '%s'", status, payload.Status)
	}
}

func TestSystemError(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	code := "ERR001"
	message := "Something went wrong"
	retryable := true

	env := b.SystemError(code, message, retryable)

	if env.Domain != DomainSystem {
		t.Errorf("expected domain '%s', got '%s'", DomainSystem, env.Domain)
	}
	if env.Type != TypeError {
		t.Errorf("expected type '%s', got '%s'", TypeError, env.Type)
	}

	payload, ok := env.Payload.(*ErrorPayload)
	if !ok {
		t.Fatalf("expected ErrorPayload type")
	}
	if payload.Code != code {
		t.Errorf("expected code '%s', got '%s'", code, payload.Code)
	}
	if payload.Message != message {
		t.Errorf("expected message '%s', got '%s'", message, payload.Message)
	}
	if payload.Retryable != retryable {
		t.Errorf("expected retryable %v, got %v", retryable, payload.Retryable)
	}
}

func TestSystemStatus(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	state := "running"

	env := b.SystemStatus(state)

	if env.Domain != DomainSystem {
		t.Errorf("expected domain '%s', got '%s'", DomainSystem, env.Domain)
	}
	if env.Type != TypeStatus {
		t.Errorf("expected type '%s', got '%s'", TypeStatus, env.Type)
	}

	payload, ok := env.Payload.(*StatusPayload)
	if !ok {
		t.Fatalf("expected StatusPayload type")
	}
	if payload.State != state {
		t.Errorf("expected state '%s', got '%s'", state, payload.State)
	}
}

func TestSystemDone(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	meta := &Meta{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		LatencyMs:        1234,
		ModelName:        "gpt-4",
		FinishReason:     "stop",
	}

	env := b.SystemDone(meta)

	if env.Domain != DomainSystem {
		t.Errorf("expected domain '%s', got '%s'", DomainSystem, env.Domain)
	}
	if env.Type != TypeDone {
		t.Errorf("expected type '%s', got '%s'", TypeDone, env.Type)
	}
	if env.Meta == nil {
		t.Fatal("expected meta to be set")
	}
	if env.Meta.TotalTokens != 150 {
		t.Errorf("expected total tokens 150, got %d", env.Meta.TotalTokens)
	}
}

func TestBusinessError(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	code := "BIZ001"
	message := "Business error"

	env := b.BusinessError(code, message)

	if env.Domain != DomainBusiness {
		t.Errorf("expected domain '%s', got '%s'", DomainBusiness, env.Domain)
	}
	if env.Type != TypeBusinessError {
		t.Errorf("expected type '%s', got '%s'", TypeBusinessError, env.Type)
	}

	payload, ok := env.Payload.(*ErrorPayload)
	if !ok {
		t.Fatalf("expected ErrorPayload type")
	}
	if payload.Code != code {
		t.Errorf("expected code '%s', got '%s'", code, payload.Code)
	}
}

func TestWorkflowStatus(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	nodeID := "node-1"
	state := "running"
	reason := "executing llm"

	env := b.WorkflowStatus(nodeID, state, reason)

	if env.Domain != DomainWorkflow {
		t.Errorf("expected domain '%s', got '%s'", DomainWorkflow, env.Domain)
	}
	if env.Type != TypeWorkflowStatus {
		t.Errorf("expected type '%s', got '%s'", TypeWorkflowStatus, env.Type)
	}

	payload, ok := env.Payload.(*WorkflowStatusPayload)
	if !ok {
		t.Fatalf("expected WorkflowStatusPayload type")
	}
	if payload.NodeID != nodeID {
		t.Errorf("expected nodeID '%s', got '%s'", nodeID, payload.NodeID)
	}
	if payload.State != state {
		t.Errorf("expected state '%s', got '%s'", state, payload.State)
	}
	if payload.Reason != reason {
		t.Errorf("expected reason '%s', got '%s'", reason, payload.Reason)
	}
}

func TestEnvelopeToJSON(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	env := b.LLMMessageDelta("Hello")

	jsonStr, err := env.ToJSON()
	if err != nil {
		t.Fatalf("failed to convert to JSON: %v", err)
	}

	// Parse and verify
	var parsed Envelope
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed.Protocol != ProtocolName {
		t.Errorf("expected protocol '%s', got '%s'", ProtocolName, parsed.Protocol)
	}
	if parsed.Domain != DomainLLM {
		t.Errorf("expected domain '%s', got '%s'", DomainLLM, parsed.Domain)
	}
}

func TestEnvelopeToSSE(t *testing.T) {
	b := NewBuilder("conv-1", "msg-1")
	env := b.LLMMessageDelta("Hello")

	sseStr, err := env.ToSSE()
	if err != nil {
		t.Fatalf("failed to convert to SSE: %v", err)
	}

	// Verify SSE format
	if !strings.HasPrefix(sseStr, "event: message\n") {
		t.Error("SSE should start with 'event: message\\n'")
	}
	if !strings.Contains(sseStr, "data: ") {
		t.Error("SSE should contain 'data: '")
	}
	if !strings.HasSuffix(sseStr, "\n\n") {
		t.Error("SSE should end with '\\n\\n'")
	}
}

func BenchmarkLLMMessageDelta(b *testing.B) {
	builder := NewBuilder("conv-1", "msg-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Reset()
		_ = builder.LLMMessageDelta("Hello, World!")
	}
}

func BenchmarkEnvelopeToJSON(b *testing.B) {
	builder := NewBuilder("conv-1", "msg-1")
	env := builder.LLMMessageDelta("Hello, World!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = env.ToJSON()
	}
}

func BenchmarkEnvelopeToSSE(b *testing.B) {
	builder := NewBuilder("conv-1", "msg-1")
	env := builder.LLMMessageDelta("Hello, World!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = env.ToSSE()
	}
}
