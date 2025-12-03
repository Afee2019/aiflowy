package tech.aiflowy.ai.provider;

import com.agentsflex.core.document.Document;
import com.agentsflex.core.util.Maps;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import dev.tinyflow.core.chain.Chain;
import dev.tinyflow.core.node.SearchEngineNode;
import dev.tinyflow.core.searchengine.BaseSearchEngine;

import java.util.*;

public class BochaaiSearchEngine extends BaseSearchEngine {

    private static final String DEFAULT_API_URL = "https://api.bochaai.com/v1/ai-search";

    public BochaaiSearchEngine() {
        this.setApiUrl(DEFAULT_API_URL);
    }

    @Override
    public List<Map<String, Object>> search(String keyword, int limit, SearchEngineNode searchEngineNode, Chain chain) {
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + this.apiKey);
        headers.put("Content-Type", "application/json");
        String jsonString = Maps.of("query", keyword).set("summary", true).set("freshness", "noLimit").set("count", limit).set("stream", false).toJSON();
//        String res = this.httpClient.post(this.apiUrl, headers, jsonString);
//        JSONObject object = JSON.parseObject(res);
//        if (200 == object.getInteger("code")) {
//            JSONArray messages = object.getJSONArray("messages");
//            List<Map<String, Object>> list = new ArrayList<>();
//            for (Object message : messages) {
//                Document document = getDocument((JSONObject) message);
//                list.add((Map<String, Object>) document);
//            }
//            return list;
//        }
        return Collections.emptyList();
    }

    private Document getDocument(JSONObject message) {
        // source|answer
        String type = message.getString("type");
        // webpage|image|video|text
        String contentType = message.getString("content_type");
        // content_type是text就是文本，其余的是json string
        String content = message.getString("content");
        Document document = new Document();
        document.setId(contentType);
        document.setTitle(type);
        document.setContent(content);
        return document;
    }
}
