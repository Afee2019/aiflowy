package tech.aiflowy.admin.controller.ai;

import tech.aiflowy.ai.entity.AiWorkflowCategory;
import tech.aiflowy.ai.service.AiWorkflowCategoryService;
import tech.aiflowy.common.annotation.UsePermission;
import tech.aiflowy.common.web.controller.BaseCurdController;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 *  控制层。
 *
 * @author ArkLight
 * @since 2025-12-11
 */
@RestController
@RequestMapping("/api/v1/aiWorkflowCategory")
@UsePermission(moduleName = "/api/v1/aiWorkflow")
public class AiWorkflowCategoryController extends BaseCurdController<AiWorkflowCategoryService, AiWorkflowCategory> {

    public AiWorkflowCategoryController(AiWorkflowCategoryService service) {
        super(service);
    }

}