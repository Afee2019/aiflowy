
package tech.aiflowy.ai.controller;

import cn.dev33.satoken.annotation.SaIgnore;
import com.agentsflex.core.llm.ChatOptions;
import com.agentsflex.core.llm.Llm;
import com.agentsflex.core.llm.functions.Function;
import com.agentsflex.core.llm.functions.Parameter;
import com.agentsflex.core.llm.response.FunctionCaller;
import com.agentsflex.core.message.FunctionCall;
import com.agentsflex.core.message.Message;
import com.agentsflex.core.message.ToolMessage;
import com.agentsflex.core.prompt.Prompt;
import com.agentsflex.core.prompt.TextPrompt;
import com.agentsflex.core.prompt.ToolPrompt;
import com.agentsflex.core.util.StringUtil;
import org.slf4j.Logger;
import org.springframework.stereotype.Controller;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import tech.aiflowy.ai.entity.AiBot;
import tech.aiflowy.ai.entity.AiBotKnowledge;
import tech.aiflowy.ai.entity.AiLlm;
import tech.aiflowy.ai.entity.AiPluginTool;
import tech.aiflowy.ai.entity.AiWorkflow;
import tech.aiflowy.ai.entity.openAi.error.OpenAiErrorResponse;
import tech.aiflowy.ai.entity.openAi.request.ChatMessage;
import tech.aiflowy.ai.entity.openAi.request.OpenAiChatRequest;
import tech.aiflowy.ai.entity.openAi.response.UnifiedLlmResponse;
import tech.aiflowy.ai.service.*;
import tech.aiflowy.common.ai.MySseEmitter;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import tech.aiflowy.ai.entity.AiKnowledge;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.ai.entity.AiBotPlugins;
import tech.aiflowy.ai.entity.AiBotWorkflow;
import com.agentsflex.core.llm.ChatContext;
import com.agentsflex.core.llm.StreamResponseListener;
import com.agentsflex.core.llm.response.AiMessageResponse;
import tech.aiflowy.common.util.Maps;
import tech.aiflowy.ai.entity.openAi.response.PlatformType;
import tech.aiflowy.ai.entity.openAi.response.ResponseConverter;
import com.alibaba.fastjson2.JSON;
import org.slf4j.LoggerFactory;

/**
 * 兼容 openAi api 的，调用 bot 的控制器
 */
@RestController
@SaIgnore
public class CompatibleChatController {

    private final Logger logger = LoggerFactory.getLogger(CompatibleChatController.class);

    @Resource
    private AiBotApiKeyService aiBotApiKeyService;

    @Resource
    private AiBotService aiBotService;

    @Resource
    private AiBotLlmService aiBotLlmService;

    @Resource
    private AiLlmService aiLlmService;

    @Resource
    private AiBotPluginsService aiBotPluginsService;

    @Resource
    private AiPluginToolService aiPluginToolService;

    @Resource
    private AiBotKnowledgeService aiBotKnowledgeService;
    ;

    @Resource
    private AiKnowledgeService aiKnowledgeService;

    @Resource
    private AiBotWorkflowService aiBotWorkflowService;

    @Resource
    private AiWorkflowService aiWorkflowService;

    @PostMapping("/v1/chat/completions")
    public Object chat(@RequestBody
                       OpenAiChatRequest params, HttpServletRequest request, HttpServletResponse response) {

        // 校验 apiKey
        String authorization = request.getHeader("Authorization");

        if (!StringUtils.hasLength(authorization)) {
            return new OpenAiErrorResponse("Invalid token", "invalid_request_error", "authorization", "401");
        }

        String apiKey = authorization.replace("Bearer ", "");
        if (!StringUtils.hasLength(apiKey)) {
            return new OpenAiErrorResponse("Invalid token", "invalid_request_error", "authorization", "401");
        }

        // 校验 messages
        List<Object> messages = params.getMessages();
        if (messages == null || messages.isEmpty()) {
            return new OpenAiErrorResponse("Bad Request", "messages can not be null or empty", "messages", "400");
        }

        BigInteger botId = null;
        try {
            botId = aiBotApiKeyService.decryptApiKey(apiKey);
        } catch (Exception e) {
            return new OpenAiErrorResponse("Invalid token", "invalid_request_error", "authorization", "401");
        }

        AiBot aiBot = aiBotService.getById(botId);
        if (aiBot == null) {
            return new OpenAiErrorResponse("Bot Not Found", "resource_not_found_error", null, "404");
        }

        // 校验 llm
        AiLlm aiLlm = aiLlmService.getById(aiBot.getLlmId());
        if (aiLlm == null) {
            return new OpenAiErrorResponse("Llm Not Found", "resource_not_found_error", null, "404");
        }

        Boolean stream = params.getStream() != null ? params.getStream() : true;

        ChatOptions chatOptions = buildChatOptions(params, aiLlm);

        List<Function> functionList = buildFunctions(aiBot, chatOptions);

        if (stream) {
            return handleStreamChat(aiLlm, chatOptions, functionList, response);
        } else {
            return handleNotStreamChat(aiLlm, chatOptions, functionList, response);
        }

    }

    private Object handleNotStreamChat(AiLlm aiLlm, ChatOptions chatOptions, List<Function> functionList,
                                       HttpServletResponse response) {

        response.setContentType("application/json;charset=utf-8");

        ServletRequestAttributes sra = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        RequestContextHolder.setRequestAttributes(sra, true);

        PlatformType platformType = PlatformType.getByBrand(aiLlm.getBrand());

        Llm llm = aiLlm.toLlm();
        chatOptions.setEnableThinking(false);
        AiMessageResponse aiResponse = llm.chat(new TextPrompt(""), chatOptions);

        if (aiResponse.isFunctionCall()){

            buildToolCallResult(aiResponse,functionList,chatOptions);
            aiResponse = llm.chat(new TextPrompt(""), chatOptions);

        }



        UnifiedLlmResponse convertResponse = ResponseConverter.handleResponse(aiResponse.getResponse(),
                platformType);
        String json = JSON.toJSONString(convertResponse);

        logger.info("大模型回复：{}", json);

        if (json.equalsIgnoreCase("{}")){
            return aiResponse.getResponse();
        }


        return convertResponse;





    }

    private SseEmitter handleStreamChat(AiLlm aiLlm, ChatOptions chatOptions, List<Function> functionList,
                                        HttpServletResponse response) {

        response.setContentType("text/event-stream");
        response.setCharacterEncoding("UTF-8");

        ServletRequestAttributes sra = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        RequestContextHolder.setRequestAttributes(sra, true);

        MySseEmitter mySseEmitter = new MySseEmitter(1000 * 60 * 300L);

        PlatformType platformType = PlatformType.getByBrand(aiLlm.getBrand());

        Boolean[] needClose = new Boolean[]{true};

        Llm llm = aiLlm.toLlm();
        llm.chatStream("", new StreamResponseListener() {
            @Override
            public void onMessage(ChatContext chatContext, AiMessageResponse aiMessageResponse) {

                if (aiMessageResponse.isFunctionCall()) {
                    needClose[0] = false;
                    buildToolCallResult(aiMessageResponse,functionList,chatOptions);
                    llm.chatStream("", new StreamResponseListener() {
                        @Override
                        public void onMessage(ChatContext context, AiMessageResponse response) {
                            UnifiedLlmResponse convertResponse = ResponseConverter.handleResponse(aiMessageResponse
                                            .getResponse(),
                                    platformType);
                            String json = JSON.toJSONString(convertResponse);
                            logger.info("大模型 function calling 回复：{}", json);
                            mySseEmitter.send(json);
                        }

                        @Override
                        public void onStop(ChatContext context) {
                            mySseEmitter.complete();
                        }

                        @Override
                        public void onFailure(ChatContext context, Throwable throwable) {
                            logger.error("function fail:{}", throwable.getMessage());
                            mySseEmitter.send(JSON.toJSONString(throwable.getMessage()));
                            mySseEmitter.complete();
                        }
                    },chatOptions);


                }

                if ("[DONE]".equalsIgnoreCase(aiMessageResponse.getResponse())) {
                    return;
                }

                UnifiedLlmResponse convertResponse = ResponseConverter.handleResponse(aiMessageResponse
                                .getResponse(),
                        platformType);
                String json = JSON.toJSONString(convertResponse);
                logger.info("大模型回复：{}", json);
                mySseEmitter.send(json);

            }



            @Override
            public void onStop(ChatContext context) {
                if (needClose[0]) {
                    mySseEmitter.complete();
                }

            }

            @Override
            public void onFailure(ChatContext context, Throwable throwable) {
                logger.error("fail:{}", throwable.getMessage());
                OpenAiErrorResponse error = new OpenAiErrorResponse(throwable.getMessage(), "error", null, "500");
                mySseEmitter.send(JSON.toJSONString(error));
                mySseEmitter.complete();
            }

        }, chatOptions);

        return mySseEmitter;

    }

    /**
     * 构建工具调用后的 chatOptions
     * @param aiMessageResponse 大模型响应
     * @param functionList 工具列表
     * @param chatOptions chat配置
     */
    private void buildToolCallResult(AiMessageResponse aiMessageResponse,List<Function> functionList,ChatOptions chatOptions) {
        List<FunctionCall> calls = aiMessageResponse.getMessage().getCalls();
        logger.info("isFunctionCall:{}", calls);

        List<FunctionCaller> functionCallers = new ArrayList<>(calls.size());
        for (FunctionCall call : calls) {
            Function function = functionList.stream().filter(fun -> fun.getName().equals(call.getName())).findFirst().orElse(null);
            if (function != null) {
                functionCallers.add(new FunctionCaller(function, call));
            }
        }

        List<Map<String, Object>> messages = (List<Map<String, Object>>) chatOptions.getExtra().get("messages");
        HashMap<String, Object> toolCallsMessage = new HashMap<>();
        toolCallsMessage.put("content", "");
        toolCallsMessage.put("role", "assistant");
        for (FunctionCaller functionCaller : functionCallers) {

            HashMap<String, Object> toolMessage = new HashMap<>();
            String callId = functionCaller.getFunctionCall().getId();


            toolMessage.put("role", "tool");

            ArrayList<Map<String, Object>> toolCalls = (ArrayList<Map<String, Object>>) toolCallsMessage.get("tool_calls");

            if (toolCalls == null) {
                toolCalls = new ArrayList<>();
            }

            if (StringUtil.hasText(callId)) {
                toolCalls.add(
                        Maps.of(
                                        "function",
                                        Maps.of(
                                                        "name", functionCaller.getFunctionCall().getName()
                                                )
                                                .set("arguments", functionCaller.getFunctionCall().getArgs())
                                )
                                .set("type", "function").set("id", callId)
                );
                toolMessage.put("tool_call_id", callId);
            } else {

                toolCalls.add(
                        Maps.of(
                                        "function",
                                        Maps.of(
                                                        "name", functionCaller.getFunctionCall().getName()
                                                )
                                                .set("arguments", JSON.toJSONString(functionCaller.getFunctionCall().getArgs()))
                                )
                                .set("type", "function")
                                .set("id", functionCaller.getFunctionCall().getName())
                );
                toolMessage.put("tool_call_id", functionCaller.getFunctionCall().getName());
            }
            toolCallsMessage.put("tool_calls", toolCalls);
            Object object = functionCaller.call();
            if (object instanceof CharSequence || object instanceof Number) {
                toolMessage.put("content", object.toString());
            } else {
                toolMessage.put("content", JSON.toJSONString(object));
            }

            messages.add(toolMessage);

        }

        messages.add(toolCallsMessage);
    }

    private ChatOptions buildChatOptions(OpenAiChatRequest params, AiLlm aiLlm) {

        ChatOptions chatOptions = params.buildChatOptions(aiLlm);

        return chatOptions;
    }

    private void buildFunctionJsonArray(List<Map<String, Object>> functionsJsonArray, List<Function> functions) {
        for (Function function : functions) {
            Map<String, Object> functionRoot = new HashMap<>();
            functionRoot.put("type", "function");

            Map<String, Object> functionObj = new HashMap<>();
            functionRoot.put("function", functionObj);

            functionObj.put("name", function.getName());
            functionObj.put("description", function.getDescription());

            Map<String, Object> parametersObj = new HashMap<>();
            functionObj.put("parameters", parametersObj);
            parametersObj.put("type", "object");

            Map<String, Object> propertiesObj = new HashMap<>();
            parametersObj.put("properties", propertiesObj);

            addParameters(function.getParameters(), propertiesObj, parametersObj);

            functionsJsonArray.add(functionRoot);
        }
    }

    private void addParameters(Parameter[] parameters, Map<String, Object> propertiesObj,
                               Map<String, Object> parametersObj) {
        List<String> requiredProperties = new ArrayList<>();
        for (Parameter parameter : parameters) {
            Map<String, Object> parameterObj = new HashMap<>();
            parameterObj.put("type", parameter.getType());
            parameterObj.put("description", parameter.getDescription());
            parameterObj.put("enum", parameter.getEnums());
            if (parameter.isRequired()) {
                requiredProperties.add(parameter.getName());
            }

            List<Parameter> children = parameter.getChildren();
            if (children != null && !children.isEmpty() && "object".equalsIgnoreCase(parameter.getType())) {
                Map<String, Object> childrenObj = new HashMap<>();
                parameterObj.put("properties", childrenObj);
                addParameters(children.toArray(new Parameter[0]), childrenObj, parameterObj);
            }

            propertiesObj.put(parameter.getName(), parameterObj);
        }

        if (!requiredProperties.isEmpty()) {
            parametersObj.put("required", requiredProperties);
        }
    }

    private ArrayList<Function> buildFunctions(AiBot aiBot, ChatOptions chatOptions) {

        BigInteger botId = aiBot.getId();
        ArrayList<Function> functionList = new ArrayList<>();

        // 绑定知识库
        List<AiBotKnowledge> aiBotKnowledgeList = aiBotKnowledgeService.listByBotId(botId);
        if (aiBotKnowledgeList != null && !aiBotKnowledgeList.isEmpty()) {
            List<BigInteger> knowledgeIds = aiBotKnowledgeList.stream()
                    .map(AiBotKnowledge::getKnowledgeId)
                    .collect(Collectors.toList());
            List<AiKnowledge> aiKnowledgeList = aiKnowledgeService.listByIds(knowledgeIds);

            if (aiKnowledgeList != null && !aiKnowledgeList.isEmpty()) {
                aiKnowledgeList.forEach(aiKnowledge -> {
                    Function function = aiKnowledge.toFunction(true);
                    functionList.add(function);
                });
            }
        }

        // 绑定插件
        List<AiPluginTool> aiPluginToolList = (List<AiPluginTool>) aiPluginToolService.getPluginToolList(botId).data();
        if (aiPluginToolList != null && !aiPluginToolList.isEmpty()) {
            aiPluginToolList.forEach(aiPluginTool -> {
                Function function = aiPluginTool.toFunction();
                functionList.add(function);
            });
        }

        // 绑定工作流
        List<AiBotWorkflow> aiBotWorkflowList = aiBotWorkflowService.listByBotId(botId);
        if (aiBotWorkflowList != null && !aiBotWorkflowList.isEmpty()) {
            List<BigInteger> workflowIds = aiBotWorkflowList.stream()
                    .map(AiBotWorkflow::getWorkflowId)
                    .collect(Collectors.toList());
            List<AiWorkflow> aiWorkflowList = aiWorkflowService.listByIds(workflowIds);

            if (aiWorkflowList != null && !aiWorkflowList.isEmpty()) {
                aiWorkflowList.forEach(workflow -> {
                    Function function = workflow.toFunction(true);
                    functionList.add(function);
                });
            }
        }

        ArrayList<Map<String, Object>> functionJsonArray = new ArrayList<>();
        buildFunctionJsonArray(functionJsonArray, functionList);

        if (functionJsonArray != null && !functionJsonArray.isEmpty()) {
            Map<String, Object> extra = chatOptions.getExtra();
            if (extra != null && !extra.isEmpty()) {

                if (extra.get("payload") != null) {
                    List<Map<String, Object>> sparkFunctions = new ArrayList<>();
                    Map<String, Object> payload = (Map<String, Object>) extra.get("payload");

                    functionJsonArray.forEach(function -> {
                        Map<String, Object> functionObj = (Map<String, Object>) function.get("function");
                        sparkFunctions.add(
                                Maps.of("name", functionObj.get("name"))
                                        .set("description", functionObj.get("description"))
                                        .set("parameters", functionObj.get("parameters"))
                                        .set("required", functionObj.get("required"))
                        );

                    });
                    payload.put("functions", Maps.of("text", sparkFunctions));

                } else {
                    extra.put("tools", functionJsonArray);
                    extra.put("tool_choice", "auto");
                }

            }
        }

        return functionList;

    }

}
