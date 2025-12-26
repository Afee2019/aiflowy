package tech.aiflowy.ai.node;

import cn.hutool.core.util.StrUtil;
import com.alibaba.fastjson.JSONObject;
import dev.tinyflow.core.node.BaseNode;
import dev.tinyflow.core.parser.BaseNodeParser;

public class MakeFileNodeParser extends BaseNodeParser {

    @Override
    public BaseNode doParse(JSONObject root, JSONObject data, JSONObject tinyflow) {
        String suffix = data.getString("suffix");
        if (StrUtil.isEmpty(suffix)) {
            suffix = "docx";
        }
        return new MakeFileNode(suffix);
    }

    public String getNodeName() {
        return "make-file";
    }
}
