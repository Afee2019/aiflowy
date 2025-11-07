package tech.aiflowy.common.audio.impl.volc;

import tech.aiflowy.common.audio.core.AudioService;
import tech.aiflowy.common.audio.core.BaseAudioClient;

import java.io.InputStream;

/**
 * 火山引擎
 */
public class VolcanoEngineAudioService implements AudioService {

    @Override
    public void textToVoiceStream(BaseAudioClient client, String sessionId, String messageId, String text) {

    }

    @Override
    public String audioToText(InputStream is) {
        return "";
    }
}
