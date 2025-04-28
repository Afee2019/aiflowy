package tech.aiflowy.ai.node;

import java.io.InputStream;

public interface ReadDocService {

    String read(String fileName, InputStream is);
}
