package tech.aiflowy.common.audio.socket;

import cn.hutool.core.codec.Base64;
import cn.hutool.core.io.FileUtil;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.socket.CloseStatus;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.web.socket.handler.TextWebSocketHandler;
import tech.aiflowy.common.audio.core.AudioServiceManager;
import tech.aiflowy.common.audio.core.BaseAudioClient;
import tech.aiflowy.common.audio.listener.AudioMessageListener;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class AudioSocketHandler extends TextWebSocketHandler {

    public static final String SESSION_ID = "sessionId";
    public static final String MESSAGE_ID = "messageId";
    public static final String TYPE = "type";
    public static final String CONTENT = "content";
    public static final String DATA = "_data_";
    public static final String END = "_end_";
    public static final String ERROR = "_error_";
    public static final String START = "_start_";

    public static Map<String, SocketEntity> sessionMap = new ConcurrentHashMap<>();

    private final AudioServiceManager manager;

    private final static Logger log = LoggerFactory.getLogger(AudioSocketHandler.class);

    public AudioSocketHandler(AudioServiceManager manager) {
        this.manager = manager;
    }

    @Override
    public void afterConnectionEstablished(WebSocketSession session) throws Exception {
        String sessionId = session.getAttributes().get(SESSION_ID).toString();
        log.info("连接建立：{}", sessionId);
        BaseAudioClient client = manager.getClient();
        SocketEntity entity = new SocketEntity();
        entity.setSessionId(sessionId);
        entity.setSession(session);
        entity.setClient(client);
        sessionMap.put(sessionId, entity);
    }

    @Override
    protected void handleTextMessage(WebSocketSession session, TextMessage message) throws Exception {
        String sessionId = session.getAttributes().get(SESSION_ID).toString();
        SocketEntity entity = getSocketEntity(sessionId);
        BaseAudioClient client = entity.getClient();
        String msg = message.getPayload();

        JSONObject obj = JSON.parseObject(msg);
        String messageId = obj.getString(MESSAGE_ID);
        String type = obj.getString(TYPE);
        String content = obj.getString(CONTENT);

        if (START.equals(type)) {
            log.info("文本转语音开始：sessionId：{}，messageId：{}", sessionId, messageId);
            handleStartMessage(entity, client, messageId);
        }
        if (END.equals(type)) {
            log.info("文本转语音结束：sessionId：{}，messageId：{}", sessionId, messageId);
            handleEndMessage(entity, client, messageId);
        }
        if (DATA.equals(type)) {
            handleDataMessage(entity, client, messageId, content);
        }
    }

    private void handleStartMessage(SocketEntity entity, BaseAudioClient client, String messageId) {
        client.addListener(entity.getSessionId(), messageId, new AudioMessageListener() {
            @Override
            public void onMessageReceived(String sessionId, String messageId, byte[] message) {
                String encode = Base64.encode(message);
                AudioSocketHandler.sendJsonVoiceMessage(sessionId, messageId, DATA, encode);
            }

            @Override
            public void onFinished(String sessionId, String messageId, byte[] fullMessage) {
                String encode = Base64.encode(fullMessage);
                /*long l = System.currentTimeMillis();
                FileUtil.writeBytes(fullMessage, "D:\\system\\desktop\\"+l+".mp3");*/
                AudioSocketHandler.sendJsonVoiceMessage(sessionId, messageId, END, encode);
            }

            @Override
            public void onError(String sessionId, String messageId, Exception e) {
                AudioSocketHandler.sendJsonVoiceMessage(sessionId, messageId, ERROR, e.getLocalizedMessage());
            }
        });
        client.beforeSend(entity.getSessionId(), messageId);
        AudioSocketHandler.sendJsonVoiceMessage(entity.getSessionId(), messageId, START, "");
    }

    private void handleEndMessage(SocketEntity entity, BaseAudioClient client, String messageId) {
        client.afterSend(entity.getSessionId(), messageId);
    }

    private void handleDataMessage(SocketEntity entity, BaseAudioClient client, String messageId, String content) {
        client.send(entity.getSessionId(), messageId, content);
    }

    @Override
    public void afterConnectionClosed(WebSocketSession session, CloseStatus status) throws Exception {
        String sessionId = session.getAttributes().get(SESSION_ID).toString();
        log.info("{} -> 断开连接：{},{}", sessionId, status.getCode(), status.getReason());
        SocketEntity entity = getSocketEntity(sessionId);
        entity.getClient().close();
    }

    public static void sendJsonVoiceMessage(String sessionId, String messageId, String msgType, String content) {
        WebSocketSession session = getSocketEntity(sessionId).getSession();
        try {
            JSONObject obj = new JSONObject();
            obj.put(MESSAGE_ID, messageId);
            obj.put(TYPE, msgType);
            obj.put(CONTENT, content);
            String msg = obj.toJSONString();
            session.sendMessage(new TextMessage(msg));
        } catch (IOException e) {
            log.error("发送语音消息失败", e);
            throw new RuntimeException(e);
        }
    }

    public static SocketEntity getSocketEntity(String sessionId) {
        SocketEntity socket = sessionMap.get(sessionId);
        if (socket == null || !socket.getSession().isOpen()) {
            log.error("获取Socket失败：WebSocket 连接为空或连接已关闭");
        }
        return socket;
    }

}
