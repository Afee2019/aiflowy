package tech.aiflowy.ai.node;

import com.alibaba.fastjson.JSONObject;
import dev.tinyflow.core.node.BaseNode;
import dev.tinyflow.core.parser.BaseNodeParser;
import tech.aiflowy.common.constant.enums.EnumResourceType;

public class DownloadNodeParser extends BaseNodeParser {

    @Override
    protected BaseNode doParse(JSONObject root, JSONObject data, JSONObject tinyflow) {
        Integer resourceType = data.getInteger("resourceType");
        if (resourceType == null) {
            resourceType = EnumResourceType.OTHER.getCode();
        }
        return new DownloadNode(resourceType);
    }

    public String getNodeName() {
        return "download-node";
    }
}
