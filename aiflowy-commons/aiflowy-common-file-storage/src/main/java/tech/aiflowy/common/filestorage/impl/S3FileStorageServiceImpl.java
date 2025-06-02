package tech.aiflowy.common.filestorage.impl;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import org.springframework.web.multipart.MultipartFile;
import tech.aiflowy.common.filestorage.FileStorageService;
import tech.aiflowy.common.filestorage.StorageConfig;
import tech.aiflowy.common.filestorage.s3.S3Client;
import tech.aiflowy.common.filestorage.s3.S3StorageConfig;

import java.io.IOException;
import java.io.InputStream;


@Component("s3")
public class S3FileStorageServiceImpl implements FileStorageService {
    private static final Logger LOG = LoggerFactory.getLogger(S3FileStorageServiceImpl.class);


    private S3Client client;

    @Value("${aiflowy.storage.s3.endpoint}")
    private  String endpoint;

    @Value("${aiflowy.storage.s3.bucket-name}")
    private  String bucketName;

    @EventListener(ApplicationReadyEvent.class)
    public void init() {
        StorageConfig instance = StorageConfig.getInstance();
        if ("s3".equals(instance.getType())) {
            client = new S3Client();
        }
    }


    @Override
    public String save(MultipartFile file) {
        try {
            return client.upload(file);
        } catch (Exception e) {
            throw new RuntimeException(e.getMessage());
        }
    }

    @Override
    public InputStream readStream(String path) throws IOException {
        return client.getObjectContent(path);
    }

    @Override
    public void delete(String path) {
        String filePath = S3FileStorageServiceImpl.extractObjectPathFromUrl(path, endpoint, bucketName);
        client.delete(filePath);
    }


    public static String extractObjectPathFromUrl(String url, String endpoint, String bucketName) {
        String host = endpoint + "/" + bucketName;

        if (url.startsWith(host)) {
            String path = url.substring(host.length());

            int queryStringIndex = path.indexOf('?');
            if (queryStringIndex > 0) {
                path = path.substring(0, queryStringIndex);
            }

            if (path.startsWith("/")) {
                path = path.substring(1);
            }

            return path;
        } else {
            throw new IllegalArgumentException("URL 不属于当前 Bucket");
        }
    }
}
