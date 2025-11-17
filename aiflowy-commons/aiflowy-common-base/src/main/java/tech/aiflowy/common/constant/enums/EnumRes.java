package tech.aiflowy.common.constant.enums;

/**
 * errorCode 枚举
 */
public enum EnumRes {

    SUCCESS(0, "成功"),
    FAIL(1, "失败"),
    NO_AUTHENTICATION(401, "请重新登陆"),
    NO_AUTHORIZATION(4010, "无权操作"),
    DUPLICATE_KEY(900, "记录已存在"),
    PARAM_ERROR(400, "参数错误"),
    ;

    private final int code;
    private final String msg;

    EnumRes(int code, String msg) {
        this.code = code;
        this.msg = msg;
    }

    public int getCode() {
        return code;
    }

    public String getMsg() {
        return msg;
    }
}
