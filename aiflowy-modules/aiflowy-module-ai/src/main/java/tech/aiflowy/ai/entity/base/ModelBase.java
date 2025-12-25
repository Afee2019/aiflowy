package tech.aiflowy.ai.entity.base;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.core.handler.FastjsonTypeHandler;
import java.io.Serializable;
import java.math.BigInteger;
import java.util.Map;


public class ModelBase implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * ID
     */
    @Id(keyType = KeyType.Generator, value = "snowFlakeId", comment = "ID")
    private BigInteger id;

    /**
     * 部门ID
     */
    @Column(comment = "部门ID")
    private BigInteger deptId;

    /**
     * 租户ID
     */
    @Column(tenantId = true, comment = "租户ID")
    private BigInteger tenantId;

    /**
     * 供应商id
     */
    @Column(comment = "供应商id")
    private BigInteger providerId;

    /**
     * 标题或名称
     */
    @Column(comment = "标题或名称")
    private String title;

    /**
     * ICON
     */
    @Column(comment = "ICON")
    private String icon;

    /**
     * 描述
     */
    @Column(comment = "描述")
    private String description;

    /**
     * 大模型请求地址
     */
    @Column(comment = "大模型请求地址")
    private String endpoint;

    /**
     * 请求路径
     */
    @Column(comment = "请求路径")
    private String requestPath;

    /**
     * 大模型名称
     */
    @Column(comment = "大模型名称")
    private String modelName;

    /**
     * 大模型 API KEY
     */
    @Column(comment = "大模型 API KEY")
    private String apiKey;

    /**
     * 大模型其他属性配置
     */
    @Column(comment = "大模型其他属性配置")
    private String extraConfig;

    /**
     * 其他配置内容
     */
    @Column(typeHandler = FastjsonTypeHandler.class, comment = "其他配置内容")
    private Map<String, Object> options;

    /**
     * 分组名称
     */
    @Column(comment = "分组名称")
    private String groupName;

    /**
     * 模型类型: chat/embedding/rerank/orc..
     */
    @Column(comment = "模型类型: chat/embedding/rerank/orc..")
    private String modelType;

    /**
     * 是否使用
     */
    @Column(comment = "是否使用")
    private Boolean withUsed;

    /**
     * 是否支持推理
     */
    @Column(comment = "是否支持推理")
    private Boolean supportThinking;

    /**
     * 是否支持工具
     */
    @Column(comment = "是否支持工具")
    private Boolean supportTool;

    /**
     * 是否支持图片
     */
    @Column(comment = "是否支持图片")
    private Boolean supportImage;

    /**
     * 仅支持 base64 的图片类型
     */
    @Column(comment = "仅支持 base64 的图片类型")
    private Boolean supportImageB64Only;

    /**
     * 是否支持视频
     */
    @Column(comment = "是否支持视频")
    private Boolean supportVideo;

    /**
     * 是否支持音频
     */
    @Column(comment = "是否支持音频")
    private Boolean supportAudio;

    /**
     * 是否免费
     */
    @Column(comment = "是否免费")
    private Boolean supportFree;

    public BigInteger getId() {
        return id;
    }

    public void setId(BigInteger id) {
        this.id = id;
    }

    public BigInteger getDeptId() {
        return deptId;
    }

    public void setDeptId(BigInteger deptId) {
        this.deptId = deptId;
    }

    public BigInteger getTenantId() {
        return tenantId;
    }

    public void setTenantId(BigInteger tenantId) {
        this.tenantId = tenantId;
    }

    public BigInteger getProviderId() {
        return providerId;
    }

    public void setProviderId(BigInteger providerId) {
        this.providerId = providerId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getIcon() {
        return icon;
    }

    public void setIcon(String icon) {
        this.icon = icon;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getEndpoint() {
        return endpoint;
    }

    public void setEndpoint(String endpoint) {
        this.endpoint = endpoint;
    }

    public String getRequestPath() {
        return requestPath;
    }

    public void setRequestPath(String requestPath) {
        this.requestPath = requestPath;
    }

    public String getModelName() {
        return modelName;
    }

    public void setModelName(String modelName) {
        this.modelName = modelName;
    }

    public String getApiKey() {
        return apiKey;
    }

    public void setApiKey(String apiKey) {
        this.apiKey = apiKey;
    }

    public String getExtraConfig() {
        return extraConfig;
    }

    public void setExtraConfig(String extraConfig) {
        this.extraConfig = extraConfig;
    }

    public Map<String, Object> getOptions() {
        return options;
    }

    public void setOptions(Map<String, Object> options) {
        this.options = options;
    }

    public String getGroupName() {
        return groupName;
    }

    public void setGroupName(String groupName) {
        this.groupName = groupName;
    }

    public String getModelType() {
        return modelType;
    }

    public void setModelType(String modelType) {
        this.modelType = modelType;
    }

    public Boolean getWithUsed() {
        return withUsed;
    }

    public void setWithUsed(Boolean withUsed) {
        this.withUsed = withUsed;
    }

    public Boolean getSupportThinking() {
        return supportThinking;
    }

    public void setSupportThinking(Boolean supportThinking) {
        this.supportThinking = supportThinking;
    }

    public Boolean getSupportTool() {
        return supportTool;
    }

    public void setSupportTool(Boolean supportTool) {
        this.supportTool = supportTool;
    }

    public Boolean getSupportImage() {
        return supportImage;
    }

    public void setSupportImage(Boolean supportImage) {
        this.supportImage = supportImage;
    }

    public Boolean getSupportImageB64Only() {
        return supportImageB64Only;
    }

    public void setSupportImageB64Only(Boolean supportImageB64Only) {
        this.supportImageB64Only = supportImageB64Only;
    }

    public Boolean getSupportVideo() {
        return supportVideo;
    }

    public void setSupportVideo(Boolean supportVideo) {
        this.supportVideo = supportVideo;
    }

    public Boolean getSupportAudio() {
        return supportAudio;
    }

    public void setSupportAudio(Boolean supportAudio) {
        this.supportAudio = supportAudio;
    }

    public Boolean getSupportFree() {
        return supportFree;
    }

    public void setSupportFree(Boolean supportFree) {
        this.supportFree = supportFree;
    }

}
