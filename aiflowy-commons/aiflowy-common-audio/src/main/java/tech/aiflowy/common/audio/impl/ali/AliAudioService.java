package tech.aiflowy.common.audio.impl.ali;

import cn.hutool.core.io.IoUtil;
import cn.hutool.core.util.IdUtil;
import cn.hutool.http.HttpRequest;
import cn.hutool.http.HttpResponse;
import cn.hutool.http.HttpUtil;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import tech.aiflowy.common.audio.config.AliConfig;
import tech.aiflowy.common.audio.core.AudioService;
import tech.aiflowy.common.audio.core.BaseAudioClient;
import tech.aiflowy.common.web.exceptions.BusinessException;

import javax.annotation.Resource;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

/**
 * 阿里云
 */
@Component("aliAudioService")
public class AliAudioService implements AudioService {

    private static final Logger log = LoggerFactory.getLogger(AliAudioService.class);
    @Resource
    private AliConfig aliConfig;

    @Override
    public void textToVoiceStream(BaseAudioClient client, String sessionId, String messageId, String text) {
        client.send(sessionId, messageId, text);
    }

    @Override
    public String audioToText(InputStream is) {

        String token = aliConfig.createToken();

        /**
         * 设置HTTPS REST POST请求
         * 1.使用http协议
         * 2.语音识别服务域名：nls-gateway-cn-shanghai.aliyuncs.com
         * 3.语音识别接口请求路径：/stream/v1/FlashRecognizer
         * 4.设置必须请求参数：appkey、token、format、sample_rate
         */
        String url = "https://nls-gateway-cn-shanghai.aliyuncs.com/stream/v1/FlashRecognizer";
        String request = url;
        request = request + "?appkey=" + aliConfig.getAppKey();
        request = request + "&token=" + token;
        request = request + "&format=" + "MP3";
        request = request + "&sample_rate=" + 16000;

        /**
         * 设置HTTPS头部字段
         * 1.Content-Type：application/octet-stream
         */
        Map<String, String> headers = new HashMap<>();
        headers.put("Content-Type", "application/octet-stream");

        HttpRequest post = HttpUtil.createPost(request);
        post.headerMap(headers, true);

        byte[] bytes = IoUtil.readBytes(is);
        post.body(bytes);

        StringBuilder sb = new StringBuilder();
        try (HttpResponse execute = post.execute()) {
            String body = execute.body();
            JSONObject obj = JSON.parseObject(body);
            Integer status = obj.getInteger("status");
            String message = obj.getString("message");
            if (20000000 != status) {
                log.error("语音识别失败：{}", obj);
                throw new BusinessException(message);
            }
            JSONArray sentences = obj.getJSONObject("flash_result")
                    .getJSONArray("sentences");

            for (Object sentence : sentences) {
                JSONObject json = (JSONObject) sentence;
                String text = json.getString("text");
                sb.append(text);
            }
        }
        return sb.toString();
    }
}
