package tech.aiflowy.ai.enums;

public enum BotMessageTypeEnum {

    NORMAL(0, "普通消息"),
    REACT_THINKING(1, "ReAct思维过程"),
    TOOL_RESULT(2, "工具调用结果");

    private final int value;
    private final String description;

    BotMessageTypeEnum(int value, String description) {
        this.value = value;
        this.description = description;
    }

    public int getValue() {
        return value;
    }

    public String getDescription() {
        return description;
    }

    public static BotMessageTypeEnum fromValue(int value) {
        for (BotMessageTypeEnum type : BotMessageTypeEnum.values()) {
            if (type.value == value) {
                return type;
            }
        }
        return NORMAL;
    }


}
