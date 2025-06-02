package tech.aiflowy.common.filestorage.impl;

import org.dromara.x.file.storage.core.FileInfo;
import org.dromara.x.file.storage.core.get.RemoteFileInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnBean;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import org.springframework.web.multipart.MultipartFile;
import tech.aiflowy.common.filestorage.FileStorageService;

import javax.annotation.Resource;
import java.io.ByteArrayInputStream;
import java.io.InputStream;


@Component("xfile")
@ConditionalOnBean(org.dromara.x.file.storage.core.FileStorageService.class)
public class XFileStorageServiceImpl implements FileStorageService {

    /**
     * 指定存储平台
     */
    @Value("${aiflowy.storage.x-file-storage.platform}")
    private String platform;

    @Resource
    private org.dromara.x.file.storage.core.FileStorageService xFileStorageService;

    @EventListener(ApplicationReadyEvent.class)
    public void init() {
    }


    @Override
    public String save(MultipartFile file) {
        try {
            FileInfo upload = xFileStorageService.of(file)
                    .setPlatform(platform)
                    .upload();
            return upload.getBasePath() + upload.getPath() + upload.getFilename();
        } catch (Exception e) {
            throw new RuntimeException(e.getMessage());
        }
    }

    @Override
    public InputStream readStream(String path) {
        RemoteFileInfo remoteFileInfo = xFileStorageService.getFile()
                .setPlatform(platform)
                .setFilename(path)
                .getFile();
        byte[] bytes = xFileStorageService.download(remoteFileInfo.toFileInfo()).bytes();
        return new ByteArrayInputStream(bytes);
    }

    @Override
    public void delete(String path) {

    }
}
