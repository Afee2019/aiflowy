package tech.aiflowy.core.chat.protocol.payload;

public class StatusPayload {
    private String state;
    private String reason;

    public String getState() {
        return state;
    }

    public void setState(String state) {
        this.state = state;
    }

    public String getReason() {
        return reason;
    }

    public void setReason(String reason) {
        this.reason = reason;
    }
}