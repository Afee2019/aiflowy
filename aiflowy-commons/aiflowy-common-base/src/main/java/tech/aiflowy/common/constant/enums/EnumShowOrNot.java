package tech.aiflowy.common.constant.enums;

import tech.aiflowy.common.annotation.DictDef;

@DictDef(name = "显示或者不显示", code = "showOrNot", keyField = "code", labelField = "text")
public enum EnumShowOrNot {

    YES(1,"显示"),
    NO(0,"不显示");

    private final int code;
    private final String text;

    EnumShowOrNot(int code, String text) {
        this.code = code;
        this.text = text;
    }

    public int getCode() {
        return code;
    }

    public String getText() {
        return text;
    }
}
