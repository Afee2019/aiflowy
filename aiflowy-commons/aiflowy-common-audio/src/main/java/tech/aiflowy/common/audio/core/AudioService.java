package tech.aiflowy.common.audio.core;

import java.io.InputStream;

public interface AudioService {

    /**
     * 文字转语音 - 流式
     */
    void textToVoiceStream(BaseAudioClient client, String sessionId, String messageId, String text);

    /**
     * 语音转文字 - 同步
     *
     * @param is 输入流
     * @return 文本
     */
    String audioToText(InputStream is);
}
