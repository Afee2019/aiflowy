package tech.aiflowy.core.chat.protocol;

public class ChatEnvelope<T> {

    private String protocol = "aiflowy-ai-chat";
    private String version = "1.1";

    private ChatDomain domain;
    private ChatType type;

    private String conversationId;
    private String messageId;
    private Integer index;

    private T payload;
    private Object meta;

    public String getProtocol() {
        return protocol;
    }

    public void setProtocol(String protocol) {
        this.protocol = protocol;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public ChatDomain getDomain() {
        return domain;
    }

    public void setDomain(ChatDomain domain) {
        this.domain = domain;
    }

    public ChatType getType() {
        return type;
    }

    public void setType(ChatType type) {
        this.type = type;
    }

    public String getConversationId() {
        return conversationId;
    }

    public void setConversationId(String conversationId) {
        this.conversationId = conversationId;
    }

    public String getMessageId() {
        return messageId;
    }

    public void setMessageId(String messageId) {
        this.messageId = messageId;
    }

    public Integer getIndex() {
        return index;
    }

    public void setIndex(Integer index) {
        this.index = index;
    }

    public T getPayload() {
        return payload;
    }

    public void setPayload(T payload) {
        this.payload = payload;
    }

    public Object getMeta() {
        return meta;
    }

    public void setMeta(Object meta) {
        this.meta = meta;
    }
}
