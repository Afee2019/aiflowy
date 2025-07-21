package tech.aiflowy.ai.entity.openAi.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.math.BigInteger;

/**
 * 统一大模型响应处理类
 * 支持平台：DeepSeek, OpenAI, 阿里百炼, 火山引擎, 百度千帆, Gitee AI, 硅基流动, Ollama
 * 星火平台使用单独的处理逻辑
 */
public class UnifiedLlmResponse {


    /**
     * 响应ID，唯一标识符
     */
    private String id;

    /**
     * 对象类型，通常为 "chat.completion"
     */
    private String object;

    /**
     * 创建时间戳
     */
    private Long created;

    /**
     * 使用的模型名称
     */
    private String model;

    /**
     * 服务等级，某些平台使用
     */
    @JsonProperty("service_tier")
    private String serviceTier;

    /**
     * 选择项列表，包含生成的内容
     */
    private List<Choice> choices;

    /**
     * 使用统计信息
     */
    private Usage usage;

    /**
     * 系统指纹，某些平台使用
     */
    @JsonProperty("system_fingerprint")
    private String systemFingerprint;

    /**
     * 选择项类
     */
    public static class Choice {
        /**
         * 选择项索引
         */
        private Integer index;

        /**
         * 生成的消息内容
         */
        private Message message;

        /**
         * 流式响应的增量内容
         */
        private Delta delta;

        /**
         * 日志概率信息
         */
        private Object logprobs;

        /**
         * 完成原因：stop, length, tool_calls, content_filter
         */
        @JsonProperty("finish_reason")
        private String finishReason;

        public Integer getIndex() {
            return index;
        }

        public void setIndex(Integer index) {
            this.index = index;
        }

        public Message getMessage() {
            return message;
        }

        public void setMessage(Message message) {
            this.message = message;
        }

        public Delta getDelta() {
            return delta;
        }

        public void setDelta(Delta delta) {
            this.delta = delta;
        }

        public Object getLogprobs() {
            return logprobs;
        }

        public void setLogprobs(Object logprobs) {
            this.logprobs = logprobs;
        }

        public String getFinishReason() {
            return finishReason;
        }

        public void setFinishReason(String finishReason) {
            this.finishReason = finishReason;
        }
    }

    /**
     * 消息类
     */
    public static class Message {
        /**
         * 角色：system, user, assistant, tool
         */
        private String role;

        /**
         * 消息内容
         */
        private String content;

        /**
         * 思考过程内容（DeepSeek-R1等推理模型使用）
         */
        @JsonProperty("reasoning_content")
        private String reasoningContent;

        /**
         * 工具调用信息
         */
        @JsonProperty("tool_calls")
        private List<ToolCall> toolCalls;

        /**
         * 工具调用ID
         */
        @JsonProperty("tool_call_id")
        private String toolCallId;

        /**
         * 名称，用于function角色
         */
        private String name;

        public String getRole() {
            return role;
        }

        public void setRole(String role) {
            this.role = role;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public String getReasoningContent() {
            return reasoningContent;
        }

        public void setReasoningContent(String reasoningContent) {
            this.reasoningContent = reasoningContent;
        }

        public List<ToolCall> getToolCalls() {
            return toolCalls;
        }

        public void setToolCalls(List<ToolCall> toolCalls) {
            this.toolCalls = toolCalls;
        }

        public String getToolCallId() {
            return toolCallId;
        }

        public void setToolCallId(String toolCallId) {
            this.toolCallId = toolCallId;
        }

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }
    }

    /**
     * 增量内容类（用于流式响应）
     */
    public static class Delta {
        /**
         * 角色
         */
        private String role;

        /**
         * 内容
         */
        private String content;

        /**
         * 思考过程内容
         */
        @JsonProperty("reasoning_content")
        private String reasoningContent;

        /**
         * 工具调用信息
         */
        @JsonProperty("tool_calls")
        private List<ToolCall> toolCalls;

        public String getRole() {
            return role;
        }

        public void setRole(String role) {
            this.role = role;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public String getReasoningContent() {
            return reasoningContent;
        }

        public void setReasoningContent(String reasoningContent) {
            this.reasoningContent = reasoningContent;
        }

        public List<ToolCall> getToolCalls() {
            return toolCalls;
        }

        public void setToolCalls(List<ToolCall> toolCalls) {
            this.toolCalls = toolCalls;
        }
    }

    /**
     * 工具调用类
     */
    public static class ToolCall {
        /**
         * 工具调用ID
         */
        private String id;

        /**
         * 类型，通常为 "function"
         */
        private String type;

        /**
         * 函数信息
         */
        private Function function;

        public static class Function {
            /**
             * 函数名称
             */
            private String name;

            /**
             * 函数参数，JSON字符串格式
             */
            private String arguments;

            public String getName() {
                return name;
            }

            public void setName(String name) {
                this.name = name;
            }

            public String getArguments() {
                return arguments;
            }

            public void setArguments(String arguments) {
                this.arguments = arguments;
            }
        }

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public String getType() {
            return type;
        }

        public void setType(String type) {
            this.type = type;
        }

        public Function getFunction() {
            return function;
        }

        public void setFunction(Function function) {
            this.function = function;
        }
    }

    /**
     * 使用统计类
     */
    public static class Usage {
        /**
         * 提示词token数量
         */
        @JsonProperty("prompt_tokens")
        private Integer promptTokens;

        /**
         * 完成内容token数量
         */
        @JsonProperty("completion_tokens")
        private Integer completionTokens;

        /**
         * 总token数量
         */
        @JsonProperty("total_tokens")
        private Integer totalTokens;

        /**
         * 缓存token数量（某些平台支持）
         */
        @JsonProperty("cached_tokens")
        private Integer cachedTokens;

        /**
         * 推理token数量（推理模型使用）
         */
        @JsonProperty("reasoning_tokens")
        private Integer reasoningTokens;

        /**
         * 提示词token详情
         */
        @JsonProperty("prompt_tokens_details")
        private PromptTokensDetails promptTokensDetails;

        /**
         * 完成token详情
         */
        @JsonProperty("completion_tokens_details")
        private CompletionTokensDetails completionTokensDetails;

        public static class PromptTokensDetails {
            @JsonProperty("cached_tokens")
            private Integer cachedTokens;

            public Integer getCachedTokens() {
                return cachedTokens;
            }

            public void setCachedTokens(Integer cachedTokens) {
                this.cachedTokens = cachedTokens;
            }
        }

        public static class CompletionTokensDetails {
            @JsonProperty("reasoning_tokens")
            private Integer reasoningTokens;

            public Integer getReasoningTokens() {
                return reasoningTokens;
            }

            public void setReasoningTokens(Integer reasoningTokens) {
                this.reasoningTokens = reasoningTokens;
            }
        }

        public Integer getPromptTokens() {
            return promptTokens;
        }

        public void setPromptTokens(Integer promptTokens) {
            this.promptTokens = promptTokens;
        }

        public Integer getCompletionTokens() {
            return completionTokens;
        }

        public void setCompletionTokens(Integer completionTokens) {
            this.completionTokens = completionTokens;
        }

        public Integer getTotalTokens() {
            return totalTokens;
        }

        public void setTotalTokens(Integer totalTokens) {
            this.totalTokens = totalTokens;
        }

        public Integer getCachedTokens() {
            return cachedTokens;
        }

        public void setCachedTokens(Integer cachedTokens) {
            this.cachedTokens = cachedTokens;
        }

        public Integer getReasoningTokens() {
            return reasoningTokens;
        }

        public void setReasoningTokens(Integer reasoningTokens) {
            this.reasoningTokens = reasoningTokens;
        }

        public PromptTokensDetails getPromptTokensDetails() {
            return promptTokensDetails;
        }

        public void setPromptTokensDetails(PromptTokensDetails promptTokensDetails) {
            this.promptTokensDetails = promptTokensDetails;
        }

        public CompletionTokensDetails getCompletionTokensDetails() {
            return completionTokensDetails;
        }

        public void setCompletionTokensDetails(CompletionTokensDetails completionTokensDetails) {
            this.completionTokensDetails = completionTokensDetails;
        }
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getObject() {
        return object;
    }

    public void setObject(String object) {
        this.object = object;
    }

    public Long getCreated() {
        return created;
    }

    public void setCreated(Long created) {
        this.created = created;
    }

    public String getModel() {
        return model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public String getServiceTier() {
        return serviceTier;
    }

    public void setServiceTier(String serviceTier) {
        this.serviceTier = serviceTier;
    }

    public List<Choice> getChoices() {
        return choices;
    }

    public void setChoices(List<Choice> choices) {
        this.choices = choices;
    }

    public Usage getUsage() {
        return usage;
    }

    public void setUsage(Usage usage) {
        this.usage = usage;
    }

    public String getSystemFingerprint() {
        return systemFingerprint;
    }

    public void setSystemFingerprint(String systemFingerprint) {
        this.systemFingerprint = systemFingerprint;
    }


    @Override
    public String toString() {
        return "UnifiedLlmResponse{" +
                "id='" + id + '\'' +
                ", object='" + object + '\'' +
                ", created=" + created +
                ", model='" + model + '\'' +
                ", serviceTier='" + serviceTier + '\'' +
                ", choices=" + choices +
                ", usage=" + usage +
                ", systemFingerprint='" + systemFingerprint + '\'' +
                '}';
    }
}