package tech.aiflowy.ai.entity;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONObject;
import com.agentsflex.core.llm.functions.BaseFunction;
import com.agentsflex.core.llm.functions.Function;
import com.agentsflex.core.llm.functions.Parameter;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mybatisflex.core.query.QueryWrapper;
import tech.aiflowy.ai.mapper.AiPluginMapper;
import tech.aiflowy.ai.service.AiPluginToolService;
import tech.aiflowy.common.ai.util.NestedParamConverter;
import tech.aiflowy.common.ai.util.PluginHttpClient;
import tech.aiflowy.common.ai.util.PluginParam;
import tech.aiflowy.common.ai.util.PluginParamConverter;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.util.SpringContextUtil;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class AiPluginFunction  implements Function {

    // 插件工具id
    private BigInteger pluginToolId;
    private String name;
    private String description;
    private Parameter[] parameters;

    public AiPluginFunction() {

    }

    public AiPluginFunction(AiPluginTool aiPluginTool) {
        this.name = aiPluginTool.getName();
        this.description = aiPluginTool.getDescription();
        this.pluginToolId = aiPluginTool.getId();
        this.parameters = getDefaultParameters();
    }

    public BigInteger getPluginToolId() {
        return pluginToolId;
    }

    public void setPluginToolId(BigInteger pluginToolId) {
        this.pluginToolId = pluginToolId;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setParameters(Parameter[] parameters) {
        this.parameters = parameters;
    }

    private AiPlugin getAiPlugin(BigInteger pluginId) {
        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("*")
                .from("tb_ai_plugin")
                .where("id = ?", pluginId);
        AiPluginMapper aiPluginMapper = SpringContextUtil.getBean(AiPluginMapper.class);
        AiPlugin aiPlugin1 = aiPluginMapper.selectOneByQuery(queryWrapper);
        return aiPlugin1;
    }

    private Parameter[] getDefaultParameters() {
        AiPluginToolService pluginToolService = SpringContextUtil.getBean(AiPluginToolService.class);
        QueryWrapper queryAiPluginToolWrapper = QueryWrapper.create()
                .select("*")
                .from("tb_ai_plugin_tool")
                .where("id = ? ", this.pluginToolId);
        AiPluginTool aiPluginTool = pluginToolService.getMapper().selectOneByQuery(queryAiPluginToolWrapper);
        List<Map<String, Object>> dataList = getDataList(aiPluginTool.getInputData());
        Parameter[] params = new Parameter[dataList.size()];
        for (int i = 0; i < dataList.size(); i++) {
            Map<String, Object> item = dataList.get(i);
            Parameter parameter = new Parameter();
            parameter.setName((String) item.get("name"));
            parameter.setDescription((String) item.get("description"));
            parameter.setRequired((boolean) item.get("required"));
            parameter.setType((String) item.get("type"));
            params[i] = parameter;
        }
        return params;
    }

    // 转换输入参数
    private List<Map<String, Object>> getDataList(String jsonArray){
        List<Map<String, Object>> dataList;
        if (jsonArray == null) {
            return new ArrayList<>();
        }
        try {
            dataList = new ObjectMapper().readValue(
                    jsonArray,
                    new TypeReference<List<Map<String, Object>>>(){}
            );
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
        return dataList;
    }


    @Override
    public String getName() {
        return name;
    }

    @Override
    public String getDescription() {
        return description;
    }

    @Override
    public Parameter[] getParameters() {
        return parameters;
    }

    @Override
    public Object invoke(Map<String, Object> argsMap) {
        AiPluginToolService pluginToolService = SpringContextUtil.getBean(AiPluginToolService.class);
        QueryWrapper queryAiPluginToolWrapper = QueryWrapper.create()
                .select("*")
                .from("tb_ai_plugin_tool")
                .where("id = ? ", this.pluginToolId);
        AiPluginTool aiPluginTool = pluginToolService.getMapper().selectOneByQuery(queryAiPluginToolWrapper);
        String method = aiPluginTool.getRequestMethod().toUpperCase();
        AiPlugin aiPlugin = getAiPlugin(aiPluginTool.getPluginId());

        String url;
        if (!StrUtil.isEmpty(aiPluginTool.getBasePath())) {
            url = aiPlugin.getBaseUrl()+aiPluginTool.getBasePath();
        } else {
            url = aiPlugin.getBaseUrl()+"/"+aiPluginTool.getName();
        }

        List<Map<String, Object>> headers = getDataList(aiPlugin.getHeaders());
        Map<String, Object> headersMap = new HashMap<>();
        for (Map<String, Object> header : headers) {
            headersMap.put((String) header.get("label"), header.get("value"));
        }
        List<PluginParam> params = new ArrayList<>();

        String authType = aiPlugin.getAuthType();
        if (!StrUtil.isEmpty(authType) && "apiKey".equals(aiPlugin.getAuthType())){
            if ("headers".equals(authType)){
                headersMap.put(aiPlugin.getTokenKey(), aiPlugin.getTokenValue());
            } else {
                PluginParam pluginParam = new PluginParam();
                pluginParam.setName(aiPlugin.getTokenKey());
                pluginParam.setDefaultValue(aiPlugin.getTokenValue());
                pluginParam.setEnabled(true);
                pluginParam.setRequired(true);
                pluginParam.setMethod("query");
                params.add(pluginParam);
            }
        }
        List<PluginParam> pluginParams = PluginParamConverter.convertFromJson(aiPluginTool.getInputData());
        Map<String, Object> nestedParams = NestedParamConverter.convertToNestedParamMap(pluginParams);

        // 准备存放不同位置的参数
        List<PluginParam> queryParams = new ArrayList<>();
        List<PluginParam> bodyParams = new ArrayList<>();
        List<PluginParam> headerParams = new ArrayList<>();
        List<PluginParam> pathParams = new ArrayList<>();

        // 遍历嵌套参数
        for (Map.Entry<String, Object> entry : nestedParams.entrySet()) {
            String paramName = entry.getKey();

            // 获取原始参数定义
            PluginParam originalParam = findOriginalParam(pluginParams, paramName);
            if (originalParam == null || !originalParam.isEnabled()) {
                continue;
            }

            // 创建参数副本以避免修改原始定义
            PluginParam requestParam = new PluginParam();
            requestParam.setName(originalParam.getName());
            requestParam.setDescription(originalParam.getDescription());
            requestParam.setRequired(originalParam.isRequired());
            requestParam.setType(originalParam.getType());
            requestParam.setEnabled(originalParam.isEnabled());
            requestParam.setMethod(originalParam.getMethod());

            // 优先级: argsMap值 > 参数默认值
            if (argsMap != null && argsMap.containsKey(paramName) && requestParam.getDefaultValue() != null) {
                // 使用大模型返回的值
                requestParam.setDefaultValue(argsMap.get(paramName));
            } else if (originalParam.getDefaultValue() != null) {
                // 使用参数定义的默认值
                requestParam.setDefaultValue(originalParam.getDefaultValue());
            }

            // 根据method分类参数
            switch (originalParam.getMethod().toLowerCase()) {
                case "query":
                    queryParams.add(requestParam);
                    break;
                case "body":
                    bodyParams.add(requestParam);
                    break;
                case "header":
                    headerParams.add(requestParam);
                    break;
                case "path":
                    pathParams.add(requestParam);
                    break;
            }
        }

        // 合并所有参数
        List<PluginParam> allParams = new ArrayList<>();
        allParams.addAll(pathParams);
        allParams.addAll(queryParams);
        allParams.addAll(bodyParams);
        allParams.addAll(headerParams);

        // 发送请求
        JSONObject result = PluginHttpClient.sendRequest(url, method, headersMap, allParams);
        return result;
    }

    // 辅助方法：根据参数名查找原始参数定义
    private PluginParam findOriginalParam(List<PluginParam> params, String name) {
        for (PluginParam param : params) {
            if (name.equals(param.getName())) {
                return param;
            }
        }
        return null;
    }

    private void processParamWithChildren(Map<String, Object> paramDef, Map<String, Object> argsMap, List<PluginParam> params) {
        boolean enabled = (boolean) paramDef.get("enabled");
        if (!enabled){
            return;
        }
        String paramName = (String) paramDef.get("name");
        PluginParam pluginParam = new PluginParam();
        pluginParam.setName(paramName);
        pluginParam.setDescription((String) paramDef.get("description"));
        pluginParam.setRequired((boolean) paramDef.get("required"));
        pluginParam.setType((String) paramDef.get("type"));
        pluginParam.setEnabled((boolean) paramDef.get("enabled"));
        pluginParam.setMethod((String) paramDef.get("method"));

        // 如果用户传了值，就用用户的值；否则用默认值
        if (paramDef.get("defaultValue") != null && !"".equals(paramDef.get("defaultValue"))) {
            pluginParam.setDefaultValue(paramDef.get("defaultValue"));
        } else if (argsMap != null && paramDef.get("name").equals(paramName) && paramDef.get("defaultValue") != null) {
            pluginParam.setDefaultValue(argsMap.get(paramName));
        }

        params.add(pluginParam);

        // 处理 children
        List<Map<String, Object>> children = (List<Map<String, Object>>) paramDef.get("children");
        if (children != null) {
            for (Map<String, Object> child : children) {
                processParamWithChildren(child, argsMap, params);
            }
        }
    }
}
