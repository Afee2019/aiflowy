package tech.aiflowy.ai.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Table;
import tech.aiflowy.ai.entity.base.AiPluginBase;

import java.util.List;


/**
 *  实体类。
 *
 * @author Administrator
 * @since 2025-04-25
 */
@Table("tb_ai_plugin")
public class AiPlugin extends AiPluginBase {


    public String getTitle() {
        return this.getName();
    }

}
