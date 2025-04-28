package tech.aiflowy.ai.node;

import cn.hutool.core.io.FileUtil;
import org.apache.tika.Tika;
import org.springframework.web.multipart.MultipartFile;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;

public class CustomFile implements MultipartFile {

    private final String fileName;
    private final byte[] bytes;

    public CustomFile(String fileName, byte[] bytes) {
        this.fileName = fileName;
        this.bytes = bytes;
    }

    @Override
    public String getName() {
        return fileName;
    }

    @Override
    public String getOriginalFilename() {
        return fileName;
    }

    @Override
    public String getContentType() {
        Tika tika = new Tika();
        InputStream inputStream = new ByteArrayInputStream(bytes);
        String contentType = "";
        try {
            contentType = tika.detect(inputStream);
        } catch (Exception e) {
            System.out.println("获取contentType错误");
        }
        return contentType;
    }

    @Override
    public boolean isEmpty() {
        return bytes == null || bytes.length == 0;
    }

    @Override
    public long getSize() {
        return bytes.length;
    }

    @Override
    public byte[] getBytes() throws IOException {
        return bytes;
    }

    @Override
    public InputStream getInputStream() throws IOException {
        return new ByteArrayInputStream(bytes);
    }

    @Override
    public void transferTo(File dest) throws IOException, IllegalStateException {
        if (dest.exists() && !dest.delete()) {
            throw new IOException(
                    "Destination file [" + dest.getAbsolutePath() + "] already exists and could not be deleted");
        }
        FileUtil.writeBytes(this.bytes, dest);
    }
}
