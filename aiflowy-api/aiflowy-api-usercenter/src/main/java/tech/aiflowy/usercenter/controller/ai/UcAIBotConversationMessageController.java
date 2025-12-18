package tech.aiflowy.usercenter.controller.ai;

import cn.dev33.satoken.annotation.SaIgnore;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import tech.aiflowy.ai.entity.AiBotConversationMessage;
import tech.aiflowy.ai.service.AiBotConversationMessageService;
import tech.aiflowy.ai.service.AiBotMessageService;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.satoken.util.SaTokenUtil;
import tech.aiflowy.common.web.controller.BaseCurdController;

import javax.annotation.Resource;
import java.util.List;

@RestController
@RequestMapping("/userCenter/conversation")
@SaIgnore
public class UcAIBotConversationMessageController extends BaseCurdController<AiBotConversationMessageService, AiBotConversationMessage> {

    @Resource
    private AiBotConversationMessageService conversationMessageService;

    @Resource
    private AiBotMessageService aiBotMessageService;

    public UcAIBotConversationMessageController(AiBotConversationMessageService service) {
        super(service);
    }

    @Override
    public Result<List<AiBotConversationMessage>> list(AiBotConversationMessage entity, Boolean asTree, String sortKey, String sortType) {
        entity.setAccountId(SaTokenUtil.getLoginAccount().getId());
        sortKey = "created";
        sortType = "desc";
        return super.list(entity, asTree, sortKey, sortType);
    }

    @Override
    protected Result<?> onSaveOrUpdateBefore(AiBotConversationMessage entity, boolean isSave) {
        entity.setAccountId(SaTokenUtil.getLoginAccount().getId());
        return super.onSaveOrUpdateBefore(entity, isSave);
    }
}
