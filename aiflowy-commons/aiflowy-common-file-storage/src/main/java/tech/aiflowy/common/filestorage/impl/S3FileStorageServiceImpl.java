package tech.aiflowy.common.filestorage.impl;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
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


}
