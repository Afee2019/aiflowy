package tech.aiflowy.common.filestorage.s3;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import tech.aiflowy.common.util.SpringContextUtil;

@Configuration
@ConfigurationProperties(prefix = "aiflowy.storage.s3")
public class S3StorageConfig {
    /**
     * 域名
     */
    private String endpoint;

    /**
     * 自定义域名
     */
    private String domain;

    /**
     * 前缀
     */
    private String prefix;

    /**
     * ACCESS_KEY
     */
    private String accessKey;

    /**
     * SECRET_KEY
     */
    private String secretKey;

    /**
     * 存储空间名
     */
    private String bucketName;

    /**
     * 存储区域
     */
    private String region;

    /**
     * 是否https（1=是）
     */
    private int isHttps = 1;

    /**
     * 桶权限类型(0private 1public 2custom)
     */
    private int accessPolicy;

    public String getDomain() {
        return domain;
    }

    public void setDomain(String domain) {
        this.domain = domain;
    }

    public String getPrefix() {
        return prefix;
    }

    public void setPrefix(String prefix) {
        this.prefix = prefix;
    }

    public String getRegion() {
        return region;
    }

    public void setRegion(String region) {
        this.region = region;
    }

    public int getIsHttps() {
        return isHttps;
    }

    public void setIsHttps(int isHttps) {
        this.isHttps = isHttps;
    }

    public int getAccessPolicy() {
        return accessPolicy;
    }

    public void setAccessPolicy(int accessPolicy) {
        this.accessPolicy = accessPolicy;
    }

    public String getEndpoint() {
        return endpoint;
    }

    public void setEndpoint(String endpoint) {
        this.endpoint = endpoint;
    }

    public String getAccessKey() {
        return accessKey;
    }

    public void setAccessKey(String accessKey) {
        this.accessKey = accessKey;
    }

    public String getSecretKey() {
        return secretKey;
    }

    public void setSecretKey(String secretKey) {
        this.secretKey = secretKey;
    }

    public String getBucketName() {
        return bucketName;
    }

    public void setBucketName(String bucketName) {
        this.bucketName = bucketName;
    }

    public static S3StorageConfig getInstance(){
        return SpringContextUtil.getBean(S3StorageConfig.class);
    }

}
