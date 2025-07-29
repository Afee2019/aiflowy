package tech.aiflowy.ai.service;

import tech.aiflowy.ai.entity.AiWorkflow;
import com.mybatisflex.core.service.IService;
import tech.aiflowy.common.domain.Result;

/**
 *  服务层。
 *
 * @author michael
 * @since 2024-08-23
 */
public interface AiWorkflowService extends IService<AiWorkflow> {

    /**
     * 根据别名或 id 查询详情
     */
    Result getDetail(String idOrAlias);


    AiWorkflow getByAlias(String alias);
}
