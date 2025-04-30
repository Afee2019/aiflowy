package tech.aiflowy.common.ai.util;
import cn.hutool.core.util.ObjectUtil;
import cn.hutool.http.*;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;

import java.util.*;

public class PluginHttpClient {

    private static final int TIMEOUT = 10_000;

    /**
     * 发送插件请求，支持 path/query/body/header 四种参数类型
     *
     * @param url           API 地址
     * @param method        HTTP 方法（GET, POST, PUT, DELETE）
     * @param headers       默认请求头
     * @param pluginParams  插件参数定义列表
     * @return              JSON 格式的响应结果
     */
    public static JSONObject sendRequest(String url, String method, Map<String, Object> headers, List<PluginParam> pluginParams) {
        // 替换 URL 中的路径变量
        String processedUrl = replacePathVariables(url, pluginParams);

        // 设置请求方法
        Method httpMethod = Method.valueOf(method.toUpperCase());
        HttpRequest request = HttpRequest.of(processedUrl).method(httpMethod);

        // 处理请求头
        if (ObjectUtil.isNotEmpty(headers)) {
            for (Map.Entry<String, Object> entry : headers.entrySet()) {
                request.header(entry.getKey(), entry.getValue().toString());
            }
        }

        // 分类参数
        Map<String, Object> queryParams = new HashMap<>();
        Map<String, Object> bodyParams = new HashMap<>();

        for (PluginParam param : pluginParams) {
            if (!param.isEnabled()) continue;

            String paramName = param.getName();
            Object paramValue = param.getDefaultValue();

            if ("Query".equalsIgnoreCase(param.getMethod())) {
                queryParams.put(paramName, paramValue);
            } else if ("Body".equalsIgnoreCase(param.getMethod())) {
                bodyParams.put(paramName, paramValue);
            } else if ("Header".equalsIgnoreCase(param.getMethod())) {
                request.header(paramName, paramValue.toString());
            }
        }

        // 添加 Query 参数
        if (!queryParams.isEmpty()) {
            request.form(queryParams); // GET 请求会自动转为 query string
        }

        // 添加 Body 参数
        if (!bodyParams.isEmpty() && (httpMethod == Method.POST || httpMethod == Method.PUT)) {
            request.body(JSONUtil.toJsonStr(bodyParams));
            request.header(Header.CONTENT_TYPE, ContentType.JSON.getValue());
        }

        // 执行请求
        HttpResponse response = request.timeout(TIMEOUT).execute();
        return JSONUtil.parseObj(response.body());
    }

    /**
     * 替换 URL 中的路径变量 {xxx}
     */
    private static String replacePathVariables(String url, List<PluginParam> pluginParams) {
        String result = url;

        // 提取 path 类型的参数
        Map<String, Object> pathParams = new HashMap<>();
        for (PluginParam param : pluginParams) {
            if ("path".equalsIgnoreCase(param.getMethod())) {
                pathParams.put(param.getName(), param.getDefaultValue());
            }
        }

        // 替换 URL 中的路径变量
        for (Map.Entry<String, Object> entry : pathParams.entrySet()) {
            result = result.replaceAll("\\{" + entry.getKey() + "\\}", entry.getValue().toString());
        }

        return result;
    }
}