package tech.aiflowy.core.chat.protocol.payload;

import java.util.Map;

public class FormRequestPayload {
    private String formId;
    private String title;
    private String description;
    private Map<String, Object> schema; // JSON Schema
    private Map<String, String> ui;

    public String getFormId() {
        return formId;
    }

    public void setFormId(String formId) {
        this.formId = formId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Map<String, Object> getSchema() {
        return schema;
    }

    public void setSchema(Map<String, Object> schema) {
        this.schema = schema;
    }

    public Map<String, String> getUi() {
        return ui;
    }

    public void setUi(Map<String, String> ui) {
        this.ui = ui;
    }
}