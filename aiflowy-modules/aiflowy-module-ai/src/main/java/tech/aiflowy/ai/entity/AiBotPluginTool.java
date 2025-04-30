package tech.aiflowy.ai.entity;


import com.mybatisflex.annotation.RelationOneToOne;
import com.mybatisflex.annotation.Table;

@Table("tb_ai_bot_plugins")
public class AiBotPluginTool {

    @RelationOneToOne(selfField = "pluginId", targetField = "id")
    private AiPluginTool aiPluginTool;

    public AiPluginTool getAiPluginTool() {
        return aiPluginTool;
    }

    public void setAiPluginTool(AiPluginTool aiPluginTool) {
        this.aiPluginTool = aiPluginTool;
    }
}
