package tech.aiflowy.common.audio.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import tech.aiflowy.common.util.SpringContextUtil;

@Configuration
@ConfigurationProperties(prefix = "aiflowy.audio")
public class AudioConfig {

    private String type;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public static AudioConfig getInstance() {
        return SpringContextUtil.getBean(AudioConfig.class);
    }
}
