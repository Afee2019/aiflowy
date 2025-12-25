package tech.aiflowy.core.chat.protocol.websocket;

import com.alibaba.fastjson.JSON;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketSession;
import tech.aiflowy.core.chat.protocol.ChatEnvelope;

import java.io.IOException;

public class ChatWebSocketSender {

    private final WebSocketSession session;

    public ChatWebSocketSender(WebSocketSession session) {
        this.session = session;
    }

    /**
     * 发送 ChatEnvelope 消息
     */
    public void send(ChatEnvelope<?> envelope) throws IOException {
        checkOpen();
        String json = JSON.toJSONString(envelope);
        session.sendMessage(new TextMessage(json));
    }

    /**
     * 发送 error 消息
     */
    public void sendError(ChatEnvelope<?> envelope) throws IOException {
        checkOpen();
        send(envelope); // 前端可根据 envelope.type 判断 error
    }

    /**
     * 发送 done 并关闭 WebSocket
     */
    public void sendDone(ChatEnvelope<?> envelope) throws IOException {
        send(envelope);
        close();
    }

    /**
     * 关闭 WebSocket
     */
    public void close() throws IOException {
        if (session.isOpen()) {
            session.close();
        }
    }

    /**
     * 检查 session 是否仍然打开
     */
    private void checkOpen() throws IOException {
        if (!session.isOpen()) {
            throw new IOException("WebSocket session is closed");
        }
    }
}
