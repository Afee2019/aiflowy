package tech.aiflowy.ai.node;

import com.agentsflex.core.chain.Chain;
import com.agentsflex.core.chain.node.BaseNode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.aiflowy.ai.utils.DocUtil;

import java.io.ByteArrayInputStream;
import java.util.HashMap;
import java.util.Map;

public class DocNode extends BaseNode {

    private static final Logger log = LoggerFactory.getLogger(DocNode.class);
    private final ReadDocService readDocService;

    public DocNode(ReadDocService readDocService) {
        this.readDocService = readDocService;
    }

    @Override
    protected Map<String, Object> execute(Chain chain) {
        Map<String, Object> map = chain.getParameterValues(this);
        Map<String, Object> res = new HashMap<>();
        String url = map.get("fileUrl").toString();
        byte[] bytes = DocUtil.downloadFile(url);
        String docContent = readDocService.read(DocUtil.getFileNameByUrl(url), new ByteArrayInputStream(bytes));
        res.put("content", docContent);
        return res;
    }
}
