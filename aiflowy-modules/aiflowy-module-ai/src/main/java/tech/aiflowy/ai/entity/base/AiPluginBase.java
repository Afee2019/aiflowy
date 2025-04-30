package tech.aiflowy.ai.entity.base;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import java.io.Serializable;
import java.math.BigInteger;
import java.util.Date;


public class AiPluginBase implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 插件id
     */
    @Id(keyType = KeyType.Generator, value = "snowFlakeId", comment = "插件id")
    private BigInteger id;

    /**
     * 图标地址
     */
    @Column(comment = "图标地址")
    private String icon;

    /**
     * 插件名称
     */
    @Column(comment = "插件名称")
    private String name;

    /**
     * 插件名称
     */
    @Column(comment = "插件名称")
    private String description;

    /**
     * 类型
     */
    @Column(comment = "类型")
    private Integer type;

    /**
     * 基础URL
     */
    @Column(comment = "基础URL")
    private String baseUrl;

    /**
     * 请求头
     */
    @Column(comment = "请求头")
    private String headers;

    /**
     * 位置 将apiKey【存放在headers 或 query中】
     */
    @Column(comment = "位置")
    private String position;

    /**
     * 认证方式
     */
    @Column(comment = "认证方式")
    private String authType;

    /**
     * tokenKey
     */
    @Column(comment = "tokenKey")
    private String tokenKey;

    /**
     * tokenValue
     */
    @Column(comment = "tokenValue")
    private String tokenValue;


    /**
     * 创建时间
     */
    @Column(comment = "创建时间")
    private Date created;

    public BigInteger getId() {
        return id;
    }

    public void setId(BigInteger id) {
        this.id = id;
    }

    public String getIcon() {
        return icon;
    }

    public void setIcon(String icon) {
        this.icon = icon;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Integer getType() {
        return type;
    }

    public void setType(Integer type) {
        this.type = type;
    }

    public String getBaseUrl() {
        return baseUrl;
    }

    public void setBaseUrl(String baseUrl) {
        this.baseUrl = baseUrl;
    }

    public String getHeaders() {
        return headers;
    }

    public void setHeaders(String headers) {
        this.headers = headers;
    }

    public String getAuthType() {
        return authType;
    }

    public void setAuthType(String authType) {
        this.authType = authType;
    }

    public Date getCreated() {
        return created;
    }

    public void setCreated(Date created) {
        this.created = created;
    }

    public String getPosition() {
        return position;
    }

    public void setPosition(String position) {
        this.position = position;
    }

    public String getTokenKey() {
        return tokenKey;
    }

    public void setTokenKey(String tokenKey) {
        this.tokenKey = tokenKey;
    }

    public String getTokenValue() {
        return tokenValue;
    }

    public void setTokenValue(String tokenValue) {
        this.tokenValue = tokenValue;
    }
}
