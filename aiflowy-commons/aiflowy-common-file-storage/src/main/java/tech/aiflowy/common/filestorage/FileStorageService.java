package tech.aiflowy.common.filestorage;

import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.InputStream;

public interface FileStorageService {


    String save(MultipartFile file);

    InputStream readStream(String path) throws IOException;

    void delete(String path);

    /**
     * 上传文件
     * @param file 文件
     * @param prePath 存储桶和文件名中间的路径（不用加斜杠）
     * @return 文件url
     */
    default String save(MultipartFile file, String prePath){
        return "";
    }
}
