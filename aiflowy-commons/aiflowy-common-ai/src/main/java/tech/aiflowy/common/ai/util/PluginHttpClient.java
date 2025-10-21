package tech.aiflowy.common.ai.util;

import cn.hutool.core.util.ObjectUtil;
import cn.hutool.http.*;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import org.springframework.http.HttpEntity;
import org.springframework.web.multipart.MultipartFile;

import java.io.ByteArrayInputStream;
import java.util.*;
import java.util.concurrent.atomic.AtomicBoolean;

public class PluginHttpClient {

    private static final int TIMEOUT = 10_000;

    public static JSONObject sendRequest(String url, String method,
                                         Map<String, Object> headers,
                                         List<PluginParam> pluginParams) {
        // 1. 处理路径参数
        String processedUrl = replacePathVariables(url, pluginParams);

        // 2. 初始化请求
        Method httpMethod = Method.valueOf(method.toUpperCase());
        HttpRequest request = HttpRequest.of(processedUrl)
                .method(httpMethod)
                .timeout(TIMEOUT);
        // 3. 处理请求头（合并默认头和参数头）
        processHeaders(request, headers, pluginParams);

        // 4. 处理查询参数和请求体
        processQueryAndBodyParams(request, httpMethod, pluginParams);

        // 5. 执行请求
        HttpResponse response = request.execute();
        return JSONUtil.parseObj(response.body());
    }

    /**
     * 处理请求头（合并默认头和参数头）
     */
    private static void processHeaders(HttpRequest request,
                                       Map<String, Object> defaultHeaders,
                                       List<PluginParam> params) {
        // 添加默认头
        if (ObjectUtil.isNotEmpty(defaultHeaders)) {
            defaultHeaders.forEach((k, v) -> request.header(k, v.toString()));
        }

        // 添加参数中指定的头
        params.stream()
                .filter(p -> "header".equalsIgnoreCase(p.getMethod()) && p.isEnabled())
                .forEach(p -> request.header(p.getName(), p.getDefaultValue().toString()));
    }

    /**
     * 处理查询参数和请求体
     */
    /**
     * 处理查询参数和请求体（新增文件参数支持）
     */
    private static void processQueryAndBodyParams(HttpRequest request,
                                                  Method httpMethod,
                                                  List<PluginParam> params) {
        Map<String, Object> queryParams = new HashMap<>();
        Map<String, Object> bodyParams = new HashMap<>();
        // 标记是否包含文件参数
        AtomicBoolean hasMultipartFile = new AtomicBoolean(false);

        // 分类参数（同时检测是否有文件）
        params.stream()
                .filter(PluginParam::isEnabled)
                .forEach(p -> {
                    String methodType = p.getMethod().toLowerCase();
                    Object paramValue = buildNestedParamValue(p);

                    // 检测是否为文件参数（MultipartFile 类型）
                    if (paramValue instanceof org.springframework.web.multipart.MultipartFile) {
                        hasMultipartFile.set(true);
                    }

                    switch (methodType) {
                        case "query":
                            queryParams.put(p.getName(), paramValue);
                            break;
                        case "body":
                            bodyParams.put(p.getName(), paramValue);
                            break;
                    }
                });

        // 1. 设置查询参数（原有逻辑不变）
        if (!queryParams.isEmpty()) {
            request.form(queryParams);
        }

        // 2. 设置请求体（分两种情况：有文件 vs 无文件）
        if (!bodyParams.isEmpty() && (httpMethod == Method.POST || httpMethod == Method.PUT)) {
            if (hasMultipartFile.get()) {
                // 2.1 包含文件参数 → 用 multipart/form-data 格式
                processMultipartBody(request, bodyParams);
            } else {
                // 2.2 无文件参数 → 保持原有 JSON 格式
                request.body(JSONUtil.toJsonStr(bodyParams))
                        .header(Header.CONTENT_TYPE, ContentType.JSON.getValue());
            }
        }
    }

    /**
     * 递归构建嵌套参数值
     * @param param 当前参数
     * @return 如果是 Object 类型，返回 Map；否则返回 defaultValue
     */
    private static Object buildNestedParamValue(PluginParam param) {
        // 如果不是 Object 类型，直接返回默认值
        if (!"Object".equalsIgnoreCase(param.getType())) {
            return param.getDefaultValue();
        }

        // 如果是 Object 类型，递归处理子参数
        Map<String, Object> nestedParams = new HashMap<>();
        if (param.getChildren() != null) {
            param.getChildren().stream()
                    .filter(PluginParam::isEnabled)
                    .forEach(child -> {
                        Object childValue = buildNestedParamValue(child); // 递归处理子参数
                        nestedParams.put(child.getName(), childValue);
                    });
        }
        return nestedParams;
    }

    /**
     * 替换URL中的路径变量 {xxx}
     */
    private static String replacePathVariables(String url, List<PluginParam> params) {
        String result = url;

        // 收集路径参数
        Map<String, Object> pathParams = new HashMap<>();
        params.stream()
                .filter(p -> "path".equalsIgnoreCase(p.getMethod()) && p.isEnabled())
                .forEach(p -> pathParams.put(p.getName(), p.getDefaultValue()));

        // 替换变量
        for (Map.Entry<String, Object> entry : pathParams.entrySet()) {
            result = result.replaceAll("\\{" + entry.getKey() + "\\}",
                    entry.getValue().toString());
        }

        return result;
    }

    private static void processMultipartBody(HttpRequest request, Map<String, Object> bodyParams) {
        // 手动设置 Content-Type 为 multipart/form-data
        request.header(Header.CONTENT_TYPE, "multipart/form-data");
        for (Map.Entry<String, Object> entry : bodyParams.entrySet()) {
            String paramName = entry.getKey();
            Object paramValue = entry.getValue();

            if (paramValue instanceof MultipartFile) {
                MultipartFile file = (MultipartFile) paramValue;
                try {
                    request.form(paramName, file.getBytes(), file.getOriginalFilename());
                } catch (Exception e) {
                    throw new RuntimeException(String.format("文件参数处理失败：参数名=%s，文件名=%s",
                            paramName, file.getOriginalFilename()), e);
                }
            } else {
                // 处理普通参数
                String valueStr;
                if (paramValue instanceof String || paramValue instanceof Number || paramValue instanceof Boolean) {
                    valueStr = paramValue.toString();
                } else {
                    valueStr = JSONUtil.toJsonStr(paramValue);
                }
                request.form(paramName, valueStr);
            }
        }
    }

}
