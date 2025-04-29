package tech.aiflowy.ai.node;

import org.springframework.stereotype.Component;
import tech.aiflowy.ai.utils.DocUtil;

import java.io.InputStream;

@Component("defaultReader")
public class DefaultReadService implements ReadDocService {

    @Override
    public String read(String fileName, InputStream is) {
        String suffix = DocUtil.getSuffix(fileName);
        if ("pdf".equals(suffix)) {
            return DocUtil.readPdfFile(is);
        } else {
            return DocUtil.readWordFile(suffix, is);
        }
    }
}
