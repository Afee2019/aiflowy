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
        List<Map<String, Object>> dataList = getDataList(aiPluginTool.getInputData());

        // 遍历插件工具定义的参数列表
        for (Map<String, Object> paramDef : dataList) {
            String paramName = (String) paramDef.get("name");

            // 检查大模型返回的参数中是否包含当前参数
            if (argsMap != null && argsMap.containsKey(paramName)) {
                PluginParam pluginParam = new PluginParam();
                pluginParam.setName(paramName);
                pluginParam.setDescription((String) paramDef.get("description"));
                pluginParam.setRequired((boolean) paramDef.get("required"));
                pluginParam.setType((String) paramDef.get("type"));
                pluginParam.setMethod((String) paramDef.get("method"));
                pluginParam.setEnabled((boolean) paramDef.get("enabled"));
                // 使用大模型返回的值，而不是默认值
                pluginParam.setDefaultValue(argsMap.get(paramName));

                params.add(pluginParam);
            } else if ((boolean) paramDef.get("required")) {
                // 如果是必填参数但大模型没有返回，使用默认值或抛出异常
                PluginParam pluginParam = new PluginParam();
                pluginParam.setName(paramName);
                pluginParam.setDescription((String) paramDef.get("description"));
                pluginParam.setRequired(true);
                pluginParam.setMethod((String) paramDef.get("method"));
                pluginParam.setEnabled((boolean) paramDef.get("enabled"));
                pluginParam.setType((String) paramDef.get("type"));
                pluginParam.setDefaultValue(paramDef.get("defaultValue"));
                params.add(pluginParam);
            }
        }

        JSONObject result = PluginHttpClient.sendRequest(url, method, headersMap, params);
        return result;
    }
}
