package tech.aiflowy.common.audio.listener;

public interface AudioMessageListener {

    void onMessageReceived(String sessionId, String messageId, byte[] message);

    void onFinished(String sessionId, String messageId, byte[] fullMessage);

    void onError(String sessionId, String messageId, Exception e);
}
