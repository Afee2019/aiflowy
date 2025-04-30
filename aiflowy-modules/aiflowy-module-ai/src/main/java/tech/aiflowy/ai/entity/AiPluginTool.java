package tech.aiflowy.ai.entity;

import com.agentsflex.core.llm.functions.Function;
import com.mybatisflex.annotation.Table;
import tech.aiflowy.ai.entity.base.AiPluginToolBase;


/**
 *  实体类。
 *
 * @author Administrator
 * @since 2025-04-27
 */
@Table("tb_ai_plugin_tool")
public class AiPluginTool extends AiPluginToolBase {



    public  Function toFunction() {
        return new AiPluginFunction(this);
    }
}
