package tech.aiflowy.ai.node;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import tech.aiflowy.common.util.SpringContextUtil;

@Component
public class ReaderManager {

    @Value("${node.reader}")
    private String reader;

    public ReadDocService getReader() {
        return SpringContextUtil.getBean(reader);
    }
}
