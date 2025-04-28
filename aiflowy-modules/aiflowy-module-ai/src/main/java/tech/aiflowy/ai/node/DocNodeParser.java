package tech.aiflowy.ai.node;

import com.agentsflex.core.chain.ChainNode;
import com.alibaba.fastjson.JSONObject;
import dev.tinyflow.core.Tinyflow;
import dev.tinyflow.core.parser.BaseNodeParser;

public class DocNodeParser extends BaseNodeParser {

    private final ReadDocService readDocService;

    public DocNodeParser(ReadDocService readDocService) {
        this.readDocService = readDocService;
    }

    @Override
    public ChainNode parse(JSONObject jsonObject, Tinyflow tinyflow) {
        JSONObject data = getData(jsonObject);
        DocNode docNode = new DocNode(readDocService);
        addParameters(docNode, data);
        addOutputDefs(docNode, data);
        return docNode;
    }

    public String getNodeName() {
        return "document-node";
    }
}
