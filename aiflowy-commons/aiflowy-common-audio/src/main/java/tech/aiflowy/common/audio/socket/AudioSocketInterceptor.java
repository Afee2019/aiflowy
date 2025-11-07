package tech.aiflowy.common.audio.socket;

import cn.dev33.satoken.stp.StpUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.server.ServerHttpRequest;
import org.springframework.http.server.ServerHttpResponse;
import org.springframework.http.server.ServletServerHttpRequest;
import org.springframework.util.StringUtils;
import org.springframework.web.socket.WebSocketHandler;
import org.springframework.web.socket.server.HandshakeInterceptor;

import java.util.Map;

public class AudioSocketInterceptor implements HandshakeInterceptor {

    private static final Logger log = LoggerFactory.getLogger(AudioSocketInterceptor.class);

    @Override
    public boolean beforeHandshake(ServerHttpRequest request, ServerHttpResponse response, WebSocketHandler wsHandler, Map<String, Object> attributes) throws Exception {

        if (request instanceof ServletServerHttpRequest) {
            ServletServerHttpRequest servletRequest = (ServletServerHttpRequest) request;
            String sessionId = servletRequest.getServletRequest().getParameter("sessionId");
            String token = servletRequest.getServletRequest().getParameter("token");
            Object loginIdByToken = StpUtil.getLoginIdByToken(token);
            if (loginIdByToken == null) {
                response.setStatusCode(org.springframework.http.HttpStatus.UNAUTHORIZED);
                return false;
            }
            if (!StringUtils.hasLength(sessionId)) {
                return false;
            }
            attributes.put("sessionId", sessionId);
            return true;
        }
        return false;
    }

    @Override
    public void afterHandshake(ServerHttpRequest request, ServerHttpResponse response, WebSocketHandler wsHandler, Exception exception) {

    }
}
