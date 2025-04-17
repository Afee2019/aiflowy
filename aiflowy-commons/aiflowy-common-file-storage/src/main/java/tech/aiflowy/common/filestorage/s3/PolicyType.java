package tech.aiflowy.common.filestorage.s3;



/**
 * minio策略配置
 *
 *
 */

public enum PolicyType {

    /**
     * 只读
     */
    READ("read-only"),

    /**
     * 只写
     */
    WRITE("write-only"),

    /**
     * 读写
     */
    READ_WRITE("read-write");

    PolicyType(String type) {
        this.type = type;
    }

    public String getType() {
        return type;
    }

    /**
     * 类型
     */
    private final String type;


}
