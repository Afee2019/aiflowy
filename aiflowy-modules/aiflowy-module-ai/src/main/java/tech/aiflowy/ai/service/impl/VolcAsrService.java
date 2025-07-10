package tech.aiflowy.ai.service.impl;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import tech.aiflowy.ai.voice.volcengine.asr.volcengine.AsrResponse;
import tech.aiflowy.ai.voice.volcengine.asr.volcengine.VolcengineAsrClient;
import tech.aiflowy.ai.service.AsrService;

import java.io.InputStream;
import java.util.Arrays;

@Service("volcAsrService")
public class VolcAsrService implements AsrService {

    private static final Logger logger = LoggerFactory.getLogger(VolcAsrService.class);

    private VolcengineAsrClient asr_client;

    @Value("${voiceInput.volcengine.app.appId}")
    private String appid ;  // 项目的 appid
    @Value("${voiceInput.volcengine.app.token}")
    private String token;  // 项目的 token
    @Value("${voiceInput.volcengine.app.cluster}")
    private String cluster;  // 请求的集群
    private final String audio_format = "raw";  // wav 或者 mp3, 根据音频类型设置


    @Override
    public String recognize(InputStream audioStream) throws Exception {


        try {
            asr_client = VolcengineAsrClient.build();
            asr_client.setAppid(appid);
            asr_client.setToken(token);
            asr_client.setCluster(cluster);
            asr_client.setFormat(audio_format);
            asr_client.setShow_utterances(true);
            asr_client.asr_sync_connect();

            byte[] b = new byte[16000];
            int len = 0;
            int count = 0;
            AsrResponse asr_response = new AsrResponse();

            while ((len = audioStream.read(b)) > 0) {
                count += 1;
                logger.info("send data pack length: {}, count {}, is_last {}", len, count, audioStream.available() == 0);
                asr_response = asr_client.asr_send(Arrays.copyOfRange(b, 0, len), audioStream.available() == 0);
            }

            StringBuilder text = new StringBuilder();

            for (AsrResponse.Result result: asr_response.getResult()) {
                text.append(result.getText());
            }

            return text.toString();

        }catch (Exception e) {
            logger.error("解析语音报错：",e);
        }finally {
            if (asr_client != null) {
                asr_client.asr_close();
            }
        }

        return "";
    }


}
