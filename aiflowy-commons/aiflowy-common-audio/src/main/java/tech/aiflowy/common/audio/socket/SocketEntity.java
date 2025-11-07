package tech.aiflowy.common.audio.socket;

import org.springframework.web.socket.WebSocketSession;
import tech.aiflowy.common.audio.core.BaseAudioClient;

public class SocketEntity {

    private String sessionId;
    private WebSocketSession session;
    private BaseAudioClient client;

    public String getSessionId() {
        return sessionId;
    }

    public void setSessionId(String sessionId) {
        this.sessionId = sessionId;
    }

    public WebSocketSession getSession() {
        return session;
    }

    public void setSession(WebSocketSession session) {
        this.session = session;
    }

    public BaseAudioClient getClient() {
        return client;
    }

    public void setClient(BaseAudioClient client) {
        this.client = client;
    }
}
