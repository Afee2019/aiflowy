package tech.aiflowy.common.audio.socket;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.socket.config.annotation.EnableWebSocket;
import org.springframework.web.socket.config.annotation.WebSocketConfigurer;
import org.springframework.web.socket.config.annotation.WebSocketHandlerRegistry;
import tech.aiflowy.common.audio.core.AudioServiceManager;

import javax.annotation.Resource;

@Configuration
@EnableWebSocket
public class AudioSocketConfig implements WebSocketConfigurer {

    @Resource
    private AudioServiceManager manager;

    @Override
    public void registerWebSocketHandlers(WebSocketHandlerRegistry registry) {

        AudioSocketHandler handler = new AudioSocketHandler(manager);
        registry.addHandler(handler, "/api/v1/aiBot/ws/audio")
                .addInterceptors(new AudioSocketInterceptor())
                .setAllowedOrigins("*");
    }
}
