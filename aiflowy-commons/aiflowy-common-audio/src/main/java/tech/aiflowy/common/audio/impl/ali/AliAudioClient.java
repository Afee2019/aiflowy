package tech.aiflowy.common.audio.impl.ali;

import cn.hutool.core.util.StrUtil;
import com.alibaba.nls.client.protocol.NlsClient;
import com.alibaba.nls.client.protocol.OutputFormatEnum;
import com.alibaba.nls.client.protocol.SampleRateEnum;
import com.alibaba.nls.client.protocol.SpeechReqProtocol;
import com.alibaba.nls.client.protocol.tts.FlowingSpeechSynthesizer;
import tech.aiflowy.common.audio.config.AliConfig;
import tech.aiflowy.common.audio.core.BaseAudioClient;
import tech.aiflowy.common.audio.listener.AudioMessageListener;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class AliAudioClient extends BaseAudioClient {

    private final AliConfig aliConfig;
    private NlsClient nlsClient;

    private final Map<String, FlowingSpeechSynthesizer> execClients = new ConcurrentHashMap<>();

    public AliAudioClient(AliConfig aliConfig) {
        this.aliConfig = aliConfig;
        initClient();
    }

    private void initClient() {
        String token = aliConfig.createToken();
        this.nlsClient = new NlsClient(token);
    }

    @Override
    public void beforeSend(String sessionId, String messageId) {
        FlowingSpeechSynthesizer execClient = getExecClient(sessionId, messageId);
        try {
            execClient.start();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void send(String sessionId, String messageId, String text) {
        FlowingSpeechSynthesizer execClient = getExecClient(sessionId, messageId);
        if (StrUtil.isNotEmpty(text)) {
            execClient.send(text);
        } else {
            execClient.getConnection().sendPing();
        }
    }

    @Override
    public void afterSend(String sessionId, String messageId) {
        FlowingSpeechSynthesizer execClient = getExecClient(sessionId, messageId);
        try {
            execClient.stop();
        } catch (Exception e) {
            log.error("【阿里云nls】客户端结束时发生异常：{}", e.getMessage());
        } finally {
            execClients.remove(sessionId + messageId);
            try {
                execClient.close();
            } catch (Exception e) {
                log.error("【阿里云nls】关闭客户端时发生异常: {}", e.getMessage());
            }
        }
    }

    @Override
    public void addListener(String sessionId, String messageId, AudioMessageListener listener) {
        FlowingSpeechSynthesizer execClient = getExecClient(sessionId, messageId);
        AliListener aliListener = (AliListener) execClient.getStreamTTSListener();
        aliListener.addListener(listener);
    }

    @Override
    public void close() {
        execClients.forEach((sessionIdMessageId, execClient) -> {
            try {
                execClient.stop();
                SpeechReqProtocol.State state = execClient.getState();
            } catch (Exception e) {
                log.error("【阿里云nls】客户端结束时发生异常：{}", e.getMessage());
            } finally {
                execClient.close();
            }
        });
        execClients.clear();
        nlsClient.shutdown();
    }

    private FlowingSpeechSynthesizer getExecClient(String sessionId, String messageId) {
        FlowingSpeechSynthesizer synthesizer = execClients.get(sessionId + messageId);
        if (synthesizer == null) {
            //创建实例，建立连接。
            try {
                synthesizer = new FlowingSpeechSynthesizer(nlsClient, new AliListener(sessionId, messageId));
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
            synthesizer.setAppKey(aliConfig.getAppKey());
            //设置返回音频的编码格式。
            synthesizer.setFormat(OutputFormatEnum.MP3);
            //设置返回音频的采样率。
            synthesizer.setSampleRate(SampleRateEnum.SAMPLE_RATE_16K);
            //发音人。注意Java SDK不支持调用超高清场景对应的发音人（例如"zhiqi"），如需调用请使用restfulAPI方式。
            synthesizer.setVoice(aliConfig.getVoice());
            //音量，范围是0~100，可选，默认50。
            synthesizer.setVolume(50);
            //语调，范围是-500~500，可选，默认是0。
            synthesizer.setPitchRate(0);
            //语速，范围是-500~500，默认是0。
            synthesizer.setSpeechRate(0);
            execClients.put(sessionId + messageId, synthesizer);
        }
        return synthesizer;
    }
}
