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
import tech.aiflowy.common.ai.util.PluginHttpClient;
import tech.aiflowy.common.ai.util.PluginParam;
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

        List<Map<String, Object>> dataList = getDataList(aiPluginTool.getInputData());

        // 遍历插件工具定义的参数列表
        for (Map<String, Object> paramDef : dataList) {
            processParamWithChildren(paramDef, argsMap, params);
        }

        JSONObject result = PluginHttpClient.sendRequest(url, method, headersMap, params);
        return result;
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
        } else if (argsMap != null && paramDef.get("name").equals(paramName)){
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
