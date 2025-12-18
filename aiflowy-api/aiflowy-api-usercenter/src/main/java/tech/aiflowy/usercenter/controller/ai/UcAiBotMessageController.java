package tech.aiflowy.usercenter.controller.ai;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import tech.aiflowy.ai.entity.AiBotMessage;
import tech.aiflowy.ai.service.AiBotMessageService;
import tech.aiflowy.common.annotation.UsePermission;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.satoken.util.SaTokenUtil;
import tech.aiflowy.common.web.controller.BaseCurdController;

import java.util.List;

/**
 * Bot 消息记录表 控制层。
 *
 * @author michael
 * @since 2024-11-04
 */
@RestController
@RequestMapping("/userCenter/aiBotMessage")
@UsePermission(moduleName = "/api/v1/aiBot")
public class UcAiBotMessageController extends BaseCurdController<AiBotMessageService, AiBotMessage> {
    private final AiBotMessageService aiBotMessageService;

    public UcAiBotMessageController(AiBotMessageService service, AiBotMessageService aiBotMessageService) {
        super(service);
        this.aiBotMessageService = aiBotMessageService;
    }

    @Override
    public Result<List<AiBotMessage>> list(AiBotMessage entity, Boolean asTree, String sortKey, String sortType) {
        entity.setAccountId(SaTokenUtil.getLoginAccount().getId());
        return super.list(entity, asTree, sortKey, sortType);
    }
}
