package tech.aiflowy.common.filestorage.s3;

import com.amazonaws.services.s3.model.CannedAccessControlList;


/**
 * 桶访问策略配置
 *
 */
public enum AccessPolicyType {

    /**
     * private
     */
    PRIVATE(0, CannedAccessControlList.Private, PolicyType.WRITE),

    /**
     * public
     */
    PUBLIC(1, CannedAccessControlList.PublicRead, PolicyType.READ),

    /**
     * custom
     */
    CUSTOM(2, CannedAccessControlList.PublicRead, PolicyType.READ);

    /**
     * 桶 权限类型
     */
    private final Integer type;

    /**
     * 文件对象 权限类型
     */
    private final CannedAccessControlList acl;

    AccessPolicyType(Integer type, CannedAccessControlList acl, PolicyType policyType) {
        this.type = type;
        this.acl = acl;
        this.policyType = policyType;
    }

    /**
     * 桶策略类型
     */
    private final PolicyType policyType;

    public Integer getType() {
        return type;
    }

    public CannedAccessControlList getAcl() {
        return acl;
    }

    public PolicyType getPolicyType() {
        return policyType;
    }

    public static AccessPolicyType getByType(Integer type) {
        for (AccessPolicyType value : values()) {
            if (value.getType().equals(type)) {
                return value;
            }
        }
        return PUBLIC;
    }

}
