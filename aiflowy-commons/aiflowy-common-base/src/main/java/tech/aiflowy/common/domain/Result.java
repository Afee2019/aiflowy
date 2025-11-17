package tech.aiflowy.common.domain;

import tech.aiflowy.common.constant.enums.EnumRes;

import java.io.Serializable;

/**
 * @author michael
 */
public class Result<T> implements Serializable {
    private static final long serialVersionUID = -8744614420977483630L;

    /**
     * 返回状态码
     * @see EnumRes
     * @mock 0
     */
    private Integer errorCode;

    /**
     * 提示消息
     *
     * @mock 成功
     */
    private String message;

    /**
     * 返回的数据
     */
    private T data;

    public static Result<Void> ok() {
        Result<Void> Result = new Result<>();
        Result.setErrorCode(EnumRes.SUCCESS.getCode());
        Result.setMessage(EnumRes.SUCCESS.getMsg());
        return Result;
    }


    public static Result<Void> ok(String msg) {
        Result<Void> Result = new Result<>();
        Result.setErrorCode(EnumRes.SUCCESS.getCode());
        Result.setMessage(msg);
        return Result;
    }

    public static <T> Result<T> ok(String msg, T data) {
        Result<T> Result = new Result<>();
        Result.setErrorCode(EnumRes.SUCCESS.getCode());
        if (msg == null || msg.isEmpty()) {
            Result.setMessage(EnumRes.SUCCESS.getMsg());
        } else {
            Result.setMessage(msg);
        }
        Result.setData(data);
        return Result;
    }

    public static Result<Void> fail(String msg) {
        Result<Void> Result = new Result<>();
        Result.setErrorCode(EnumRes.FAIL.getCode());
        Result.setMessage(msg);
        return Result;
    }

    public static Result<Void> fail(int code, String msg) {
        Result<Void> Result = new Result<>();
        Result.setErrorCode(code);
        Result.setMessage(msg);
        return Result;
    }

    public static <T> Result<T> fail(String msg, T data) {
        Result<T> Result = new Result<>();
        Result.setErrorCode(EnumRes.FAIL.getCode());
        Result.setMessage(msg);
        Result.setData(data);
        return Result;
    }

    public int getErrorCode() {
        return errorCode;
    }

    public void setErrorCode(int errorCode) {
        this.errorCode = errorCode;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public T getData() {
        return data;
    }

    public void setData(T data) {
        this.data = data;
    }
}
