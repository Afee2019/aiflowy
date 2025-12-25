package tech.aiflowy.core.chat.protocol.sse;

import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import tech.aiflowy.core.chat.protocol.ChatDomain;
import tech.aiflowy.core.chat.protocol.ChatEnvelope;
import tech.aiflowy.core.chat.protocol.ChatType;
import tech.aiflowy.core.chat.protocol.payload.ErrorPayload;

public class ChatSseUtil {

    private static final long DEFAULT_TIMEOUT = 1000 * 60 * 2; // 2分钟

    /**
     * 快速发送一次系统错误消息，并立即关闭 SSE
     *
     * @param conversationId 会话ID
     * @param message        错误提示信息
     * @param code           错误码，可选
     * @return SseEmitter
     */
    public static SseEmitter sendSystemError(String conversationId, String message, String code) {
        ChatSseEmitter emitter = new ChatSseEmitter(DEFAULT_TIMEOUT);

        ChatEnvelope<ErrorPayload> envelope = new ChatEnvelope<>();
//        envelope.setProtocol("aiflowy-ai-chat");
//        envelope.setVersion("1.1");
        envelope.setDomain(ChatDomain.SYSTEM);
        envelope.setType(ChatType.ERROR);
        envelope.setConversationId(conversationId);

        ErrorPayload payload = new ErrorPayload();
        payload.setMessage(message);
        payload.setCode(code != null ? code : "SYSTEM_ERROR");
        payload.setRetryable(false);

        envelope.setPayload(payload);

        // 发送并立即关闭
        emitter.sendAndClose(envelope);

        return emitter.getEmitter();
    }

    /**
     * 快速发送一次系统错误消息（默认错误码）
     */
    public static SseEmitter sendSystemError(String conversationId, String message) {
        return sendSystemError(conversationId, message, null);
    }


    /**
     * 快速发送任意 ChatEnvelope 消息并立即关闭 SSE
     *
     * @param conversationId 会话ID
     * @param domain         消息 domain
     * @param type           消息 type
     * @param payload        消息 payload
     * @param meta           可选 meta 信息
     * @return SseEmitter 已经发送并关闭
     */
    public static <T> SseEmitter sendAndClose(
            String conversationId,
            ChatDomain domain,
            ChatType type,
            T payload,
            Object meta
    ) {
        ChatSseEmitter emitter = new ChatSseEmitter(DEFAULT_TIMEOUT);

        ChatEnvelope<T> envelope = new ChatEnvelope<>();
//        envelope.setProtocol("aiflowy-ai-chat");
//        envelope.setVersion("1.1");
        envelope.setConversationId(conversationId);
        envelope.setDomain(domain);
        envelope.setType(type);
        envelope.setPayload(payload);
        envelope.setMeta(meta);

        // 发送并立即关闭
        emitter.sendAndClose(envelope);

        return emitter.getEmitter();
    }

    /**
     * 快速发送任意 ChatEnvelope 消息，无 meta
     */
    public static <T> SseEmitter sendAndClose(
            String conversationId,
            ChatDomain domain,
            ChatType type,
            T payload
    ) {
        return sendAndClose(conversationId, domain, type, payload, null);
    }
}
