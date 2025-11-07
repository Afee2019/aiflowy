package tech.aiflowy.common.audio.core;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.aiflowy.common.audio.listener.AudioMessageListener;

public abstract class BaseAudioClient {

    protected Logger log = LoggerFactory.getLogger(getClass());

    public abstract void beforeSend(String sessionId, String messageId);

    public abstract void send(String sessionId, String messageId, String text);

    public abstract void afterSend(String sessionId, String messageId);

    public abstract void addListener(String sessionId, String messageId, AudioMessageListener listener);

    public abstract void close();
}
