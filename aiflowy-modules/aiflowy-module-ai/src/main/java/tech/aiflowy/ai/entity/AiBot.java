package tech.aiflowy.ai.entity;

import tech.aiflowy.ai.entity.base.AiBotBase;
import com.mybatisflex.annotation.Table;

import java.util.Map;

/**
 * 实体类。
 *
 * @author michael
 * @since 2024-08-23
 */

@Table("tb_ai_bot")
public class AiBot extends AiBotBase {


    public boolean isAnonymousEnabled() {
        Map<String, Object> options = getOptions();
        if (options == null) {
            return false;
        }
        Object o = options.get("anonymousEnabled");
        return o != null && (boolean) o;
    }

}
