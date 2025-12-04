package tech.aiflowy.admin.controller.listener;

import cn.hutool.core.util.ObjectUtil;
import com.agentsflex.core.message.AiMessage;
import com.agentsflex.core.message.Message;
import com.agentsflex.core.message.ToolCall;
import com.agentsflex.core.message.ToolMessage;
import com.agentsflex.core.model.chat.ChatModel;
import com.agentsflex.core.model.chat.ChatOptions;
import com.agentsflex.core.model.chat.StreamResponseListener;
import com.agentsflex.core.model.chat.response.AiMessageResponse;
import com.agentsflex.core.model.chat.tool.ToolContext;
import com.agentsflex.core.model.client.StreamContext;
import com.agentsflex.core.prompt.MemoryPrompt;
import com.agentsflex.core.util.CollectionUtil;
import com.alibaba.fastjson.JSON;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import tech.aiflowy.common.ai.ChatSseEmitter;
import tech.aiflowy.common.util.Maps;
import tech.aiflowy.common.util.StringUtil;

import java.io.IOException;
import java.util.List;

public class ChatStreamListener implements StreamResponseListener {
    private static final Logger logger = LoggerFactory.getLogger(ChatStreamListener.class);

    private final SseEmitter emitter;

    private boolean hasFinish = false;

    private ChatOptions chatOptions;

    private ChatModel chatModel;

    private MemoryPrompt memoryPrompt;

    public ChatStreamListener() {
        this.emitter = ChatSseEmitter.create();
    }

    public ChatStreamListener(ChatOptions chatOptions, ChatModel chatModel, MemoryPrompt memoryPrompt) {
        this.emitter = ChatSseEmitter.create();
        this.chatOptions = chatOptions;
        this.chatModel = chatModel;
        this.memoryPrompt = memoryPrompt;
    }

    public SseEmitter getEmitter() {
        return this.emitter;
    }

    @Override
    public void onStart(StreamContext context) {
        logger.info("[流式对话开始]");
    }

    @Override
    public void onMessage(StreamContext streamContext, AiMessageResponse aiMessageResponse) {
        String fullText = aiMessageResponse.getMessage().getFullContent();
        String delta = aiMessageResponse.getMessage().getContent();
        try {
            if (StringUtil.hasText(delta)) {
                emitter.send(delta);
            }
        } catch (Exception e) {
            logger.error("[流式推送失败]", e);
            emitter.completeWithError(e);
        }
        System.out.print(delta);
    }

    @Override
    public void onStop(StreamContext context) {
        System.out.println("结束");
        try {
            AiMessage aiMessage = context.getAiMessage();
            List<ToolCall> toolCalls = aiMessage.getToolCalls();
            if (CollectionUtil.hasItems(aiMessage.getToolCalls())){
                chatModel.chatStream(memoryPrompt, this, chatOptions);
            } else {
                emitter.send(SseEmitter.event().name("finishEvent").data(JSON.toJSONString(Maps.of("consumeTokenInfo", aiMessage))));
                emitter.complete();
                logger.info("[流式正常结束]");
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void onFailure(StreamContext context, Throwable throwable) {
        logger.error("[流式异常结束]", throwable);
        emitter.complete();
    }

}
