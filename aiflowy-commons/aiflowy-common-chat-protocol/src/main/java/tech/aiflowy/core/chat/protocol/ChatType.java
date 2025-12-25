package tech.aiflowy.core.chat.protocol;

public enum ChatType {
    THINKING,
    MESSAGE,
    TOOL_CALL,
    TOOL_RESULT,
    STATUS,
    ERROR,
    FORM_REQUEST,
    FORM_CANCEL,
    DONE
}
