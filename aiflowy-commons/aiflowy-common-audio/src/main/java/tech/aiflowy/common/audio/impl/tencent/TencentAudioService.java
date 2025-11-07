package tech.aiflowy.common.audio.impl.tencent;

import tech.aiflowy.common.audio.core.AudioService;
import tech.aiflowy.common.audio.core.BaseAudioClient;

import java.io.InputStream;

/**
 * 腾讯云
 */
public class TencentAudioService implements AudioService {

    @Override
    public void textToVoiceStream(BaseAudioClient client, String sessionId, String messageId, String text) {

    }

    @Override
    public String audioToText(InputStream is) {
        return "";
    }
}
