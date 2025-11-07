package tech.aiflowy.common.audio.impl.ali;

import com.alibaba.fastjson.JSON;
import com.alibaba.nls.client.protocol.tts.FlowingSpeechSynthesizerListener;
import com.alibaba.nls.client.protocol.tts.FlowingSpeechSynthesizerResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.aiflowy.common.audio.listener.AudioMessageListener;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.List;

public class AliListener extends FlowingSpeechSynthesizerListener {

    private final static Logger log = LoggerFactory.getLogger(AliListener.class);

    private final String sessionId;
    private final String messageId;

    private final List<AudioMessageListener> listeners = new ArrayList<>();

    private boolean firstRecvBinary = true;
    private final ByteArrayOutputStream outputStream = new ByteArrayOutputStream();

    public AliListener(String sessionId, String messageId) {
        this.sessionId = sessionId;
        this.messageId = messageId;
    }

    public void addListener(AudioMessageListener listener) {
        listeners.add(listener);
    }

    //流入语音合成开始
    public void onSynthesisStart(FlowingSpeechSynthesizerResponse response) {
        //log.info("【阿里云nls】流入语音合成开始 ---> name: {}，status: {}", response.getName(), response.getStatus());
    }

    //服务端检测到了一句话的开始
    public void onSentenceBegin(FlowingSpeechSynthesizerResponse response) {
        //log.info("【阿里云nls】服务端检测到了一句话的开始 ---> name: {}，status: {}", response.getName(), response.getStatus());
    }

    //服务端检测到了一句话的结束，获得这句话的起止位置和所有时间戳
    public void onSentenceEnd(FlowingSpeechSynthesizerResponse response) {
        //log.info("【阿里云nls】服务端检测到了一句话的结束 ---> name: {}，status: {}，subtitles: {}", response.getName(), response.getStatus(), response.getObject("subtitles"));
    }

    //流入语音合成结束
    @Override
    public void onSynthesisComplete(FlowingSpeechSynthesizerResponse response) {
        // 调用onSynthesisComplete时，表示所有TTS数据已经接收完成，所有文本都已经合成音频并返回。
        //log.info("【阿里云nls】流入语音合成结束 ---> name: {}，status: {}", response.getName(), response.getStatus());
        for (AudioMessageListener listener : listeners) {
            listener.onFinished(sessionId, messageId, outputStream.toByteArray());
        }
    }

    //收到语音合成的语音二进制数据
    @Override
    public void onAudioData(ByteBuffer message) {
        //log.info("【阿里云nls】收到语音合成的语音二进制数据。");
        if (firstRecvBinary) {
            // 此处计算首包语音流的延迟，收到第一包语音流时，即可以进行语音播放，以提升响应速度（特别是实时交互场景下）。
            firstRecvBinary = false;
        }
        byte[] bytesArray = new byte[message.remaining()];
        message.get(bytesArray, 0, bytesArray.length);
        try {
            outputStream.write(bytesArray);
        } catch (IOException e) {
            log.error("【阿里云nls】写入二进制数据失败", e);
            throw new RuntimeException(e);
        }
        for (AudioMessageListener listener : listeners) {
            listener.onMessageReceived(sessionId, messageId, bytesArray);
        }
    }

    //收到语音合成的增量音频时间戳
    @Override
    public void onSentenceSynthesis(FlowingSpeechSynthesizerResponse response) {
        //log.info("【阿里云nls】收到语音合成的增量音频时间戳 ---> name: {}，status: {}，subtitles: {}", response.getName(), response.getStatus(), response.getObject("subtitles"));
    }

    @Override
    public void onFail(FlowingSpeechSynthesizerResponse response) {
        // task_id是调用方和服务端通信的唯一标识，当遇到问题时，需要提供此task_id以便排查。
        int status = response.getStatus();
        if (status != 40000004) {
            log.error("【阿里云nls】合成失败 ---> 会话id：{}，消息id：{} session_id: {}，task_id: {}，status: {}，status_text: {}",
                    sessionId,
                    messageId,
                    getFlowingSpeechSynthesizer().getCurrentSessionId(),
                    response.getTaskId(),
                    status,
                    response.getStatusText());
            for (AudioMessageListener listener : listeners) {
                listener.onError(sessionId, messageId, new Exception(JSON.toJSONString(response)));
            }
        }
    }

    public String getSessionId() {
        return sessionId;
    }

    public String getMessageId() {
        return messageId;
    }

}
