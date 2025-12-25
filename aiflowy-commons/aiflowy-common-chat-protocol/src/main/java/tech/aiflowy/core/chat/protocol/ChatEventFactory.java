package tech.aiflowy.core.chat.protocol;

import tech.aiflowy.core.chat.protocol.payload.MessageDeltaPayload;
import tech.aiflowy.core.chat.protocol.payload.ThinkingPayload;

public final class ChatEventFactory {

    private ChatEventFactory() {
    }

    public static ChatEnvelope<ThinkingPayload> thinking(String conversationId, String content) {
        ChatEnvelope<ThinkingPayload> e = base(conversationId);
        e.setDomain(ChatDomain.LLM);
        e.setType(ChatType.THINKING);

        ThinkingPayload p = new ThinkingPayload();
        p.setContent(content);
//        p.setVisibility("hidden");

        e.setPayload(p);
        return e;
    }

    public static ChatEnvelope<MessageDeltaPayload> messageDelta(
            String conversationId,
            String delta
    ) {
        ChatEnvelope<MessageDeltaPayload> e = base(conversationId);
        e.setDomain(ChatDomain.LLM);
        e.setType(ChatType.MESSAGE);

        MessageDeltaPayload p = new MessageDeltaPayload();
        p.setDelta(delta);

        e.setPayload(p);
        return e;
    }

    public static ChatEnvelope<Void> done(String conversationId) {
        ChatEnvelope<Void> e = base(conversationId);
        e.setDomain(ChatDomain.SYSTEM);
        e.setType(ChatType.DONE);
        e.setPayload(null);
        return e;
    }

    private static <T> ChatEnvelope<T> base(String conversationId) {
        ChatEnvelope<T> e = new ChatEnvelope<>();
        e.setConversationId(conversationId);
        return e;
    }
}

