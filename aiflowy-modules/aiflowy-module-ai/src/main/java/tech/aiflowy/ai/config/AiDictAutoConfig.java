package tech.aiflowy.ai.config;

import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.event.EventListener;
import tech.aiflowy.ai.mapper.AiWorkflowCategoryMapper;
import tech.aiflowy.ai.mapper.AiWorkflowMapper;
import tech.aiflowy.common.util.SpringContextUtil;
import tech.aiflowy.core.dict.DictManager;
import tech.aiflowy.core.dict.loader.DbDataLoader;

import javax.annotation.Resource;

@Configuration
public class AiDictAutoConfig {

    @Resource
    private AiWorkflowMapper aiWorkflowMapper;
    @Resource
    private AiWorkflowCategoryMapper aiWorkflowCategoryMapper;

    @EventListener(ApplicationReadyEvent.class)
    public void onApplicationStartup() {

        DictManager dictManager = SpringContextUtil.getBean(DictManager.class);
        dictManager.putLoader(new DbDataLoader<>("aiWorkFlow", aiWorkflowMapper, "id", "title", null, null, false));
        dictManager.putLoader(new DbDataLoader<>("aiWorkFlowCategory", aiWorkflowCategoryMapper, "id", "category_name", null, null, false));
    }
}
