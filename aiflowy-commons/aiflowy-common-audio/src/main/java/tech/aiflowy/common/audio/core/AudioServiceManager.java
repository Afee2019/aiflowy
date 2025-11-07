package tech.aiflowy.common.audio.core;

import org.springframework.scheduling.TaskScheduler;
import org.springframework.stereotype.Component;
import tech.aiflowy.common.audio.config.AliConfig;
import tech.aiflowy.common.audio.config.AudioConfig;
import tech.aiflowy.common.audio.impl.ali.AliAudioClient;
import tech.aiflowy.common.util.SpringContextUtil;

import javax.annotation.Resource;
import java.io.InputStream;

@Component
public class AudioServiceManager implements AudioService {

    @Resource
    private AliConfig aliConfig;
    @Resource(name = "taskScheduler")
    private TaskScheduler scheduler;

    @Override
    public void textToVoiceStream(BaseAudioClient client, String sessionId,String messageId, String text) {
        getService().textToVoiceStream(client, sessionId,messageId, text);
    }

    @Override
    public String audioToText(InputStream is) {
        return getService().audioToText(is);
    }

    public BaseAudioClient getClient() {
        String type = AudioConfig.getInstance().getType();
        if ("aliAudioService".equals(type)) {
            return new AliAudioClient(aliConfig);
        }
        return null;
    }

    private AudioService getService() {
        String type = AudioConfig.getInstance().getType();
        return SpringContextUtil.getBean(type);
    }

    public AliConfig getAliConfig() {
        return aliConfig;
    }

}
