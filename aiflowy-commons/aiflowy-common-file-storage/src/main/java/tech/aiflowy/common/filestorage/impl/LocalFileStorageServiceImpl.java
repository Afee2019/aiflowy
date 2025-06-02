package tech.aiflowy.common.filestorage.impl;

import cn.hutool.core.io.FileTypeUtil;
import cn.hutool.core.io.file.FileNameUtil;
import cn.hutool.core.util.StrUtil;
import tech.aiflowy.common.util.DateUtil;
import tech.aiflowy.common.filestorage.FileStorageService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import java.io.*;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.util.Date;
import java.util.UUID;


@Component("local")
public class LocalFileStorageServiceImpl implements FileStorageService {
    private static final Logger LOG = LoggerFactory.getLogger(LocalFileStorageServiceImpl.class);


    @Value("${aiflowy.storage.local.root:}")
    private String root;

    @EventListener(ApplicationReadyEvent.class)
    public void init() {
    }


    @Override
    public String save(MultipartFile file) {
        try {
            byte[] bytes = file.getBytes();
            String path = generatePath(bytes, file.getOriginalFilename());
            File target = getLocalFile(path);
            if (!target.getParentFile().exists() && !target.getParentFile().mkdirs()) {
                LOG.error("Can not mkdirs: {} ", target.getParentFile());
            }
            file.transferTo(target);
            return path;
        } catch (Exception e) {
            throw new RuntimeException(e.getMessage());
        }
    }

    @Override
    public InputStream readStream(String path) throws IOException {
        File target = getLocalFile(path);
        return Files.newInputStream(target.toPath());
    }

    @Override
    public void delete(String path) {

    }


    private File getLocalFile(String path) throws IOException {
        String root = StringUtils.hasText(this.root) ? this.root : getDefaultRoot();
        return new File(root, path);
    }


    public static String getExtName(String fileName) {
        if (!StringUtils.hasText(fileName)) {
            return null;
        }
        int index = fileName.lastIndexOf('.');
        if (index != -1 && (index + 1) < fileName.length()) {
            return fileName.substring(index + 1);
        } else {
            return null;
        }
    }

    public static String newFile(String extName) {
        return "/attachment/"
                + DateUtil.toString(new Date(), "yyyy/MM-dd") + "/"
                + UUID.randomUUID() + "." + extName;
    }

    public static String generatePath(byte[] content, String originalName) throws Exception {
        // 计算文件内容的 SHA256 哈希值
        String sha256Hex = sha256Hex(content);
        // 情况一：如果存在原始文件名，优先使用其后缀
        if (StrUtil.isNotBlank(originalName)) {
            // 提取文件后缀
            String extName = FileNameUtil.extName(originalName);
            // 如果后缀存在，返回 "哈希值.后缀"，否则返回 "哈希值"
            return "/attachment/" + (StrUtil.isBlank(extName) ? sha256Hex : sha256Hex + "." + extName);
        }
        // 情况二：如果原始文件名为空，基于文件内容推断文件类型
        return "/attachment/" + sha256Hex + '.' + FileTypeUtil.getType(new ByteArrayInputStream(content));
    }

    public static String sha256Hex(byte[] input) throws Exception {
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        byte[] hash = digest.digest(input);

        StringBuilder hexString = new StringBuilder();
        for (byte b : hash) {
            String hex = Integer.toHexString(0xff & b);
            if (hex.length() == 1) hexString.append('0');
            hexString.append(hex);
        }
        return hexString.toString(); // 返回字符串类型的哈希值
    }

    private String getDefaultRoot() throws IOException {
        ClassPathResource fileResource = new ClassPathResource("/");
        return new File(fileResource.getFile(), "/public").getAbsolutePath();
    }
}
