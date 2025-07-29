package tech.aiflowy.ai.service;

import tech.aiflowy.ai.entity.AiBot;
import com.mybatisflex.core.service.IService;
import tech.aiflowy.common.domain.Result;

/**
 *  服务层。
 *
 * @author michael
 * @since 2024-08-23
 */
public interface AiBotService extends IService<AiBot> {

    Result getDetail(String id);

    void updateBotLlmId(AiBot aiBot);

    AiBot getByAlias(String alias);
}
