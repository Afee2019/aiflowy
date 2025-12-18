package tech.aiflowy.usercenter.controller.ai;

import cn.hutool.core.collection.CollectionUtil;
import com.mybatisflex.core.query.QueryWrapper;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import tech.aiflowy.ai.entity.AiBot;
import tech.aiflowy.ai.entity.AiBotRecentlyUsed;
import tech.aiflowy.ai.service.AiBotRecentlyUsedService;
import tech.aiflowy.ai.service.AiBotService;
import tech.aiflowy.common.annotation.UsePermission;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.entity.LoginAccount;
import tech.aiflowy.common.satoken.util.SaTokenUtil;
import tech.aiflowy.common.web.controller.BaseCurdController;

import javax.annotation.Resource;
import java.io.Serializable;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Date;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 最近使用 控制层。
 *
 * @author ArkLight
 * @since 2025-12-18
 */
@RestController
@RequestMapping("/userCenter/aiBotRecentlyUsed")
@UsePermission(moduleName = "/api/v1/aiBot")
public class UcBotRecentlyUsedController extends BaseCurdController<AiBotRecentlyUsedService, AiBotRecentlyUsed> {

    @Resource
    private AiBotService aiBotService;

    public UcBotRecentlyUsedController(AiBotRecentlyUsedService service) {
        super(service);
    }

    @GetMapping("/getRecentlyBot")
    public Result<List<AiBot>> getRecentlyBot() {
        LoginAccount account = SaTokenUtil.getLoginAccount();
        QueryWrapper w = QueryWrapper.create();
        w.eq(AiBotRecentlyUsed::getCreatedBy,account.getId());
        List<AiBotRecentlyUsed> list = service.list(w);
        if (CollectionUtil.isNotEmpty(list)) {
            List<BigInteger> botIds = list.stream().map(AiBotRecentlyUsed::getBotId).collect(Collectors.toList());
            QueryWrapper botQw = QueryWrapper.create();
            botQw.in(AiBot::getId,botIds);
            return Result.ok(aiBotService.list(botQw));
        }
        return Result.ok(new ArrayList<>());
    }

    @GetMapping("/removeByBotId")
    public Result<Void> removeByBotId(BigInteger botId) {
        QueryWrapper w = QueryWrapper.create();
        w.eq(AiBotRecentlyUsed::getBotId,botId);
        w.eq(AiBotRecentlyUsed::getCreatedBy,SaTokenUtil.getLoginAccount().getId());
        service.remove(w);
        return Result.ok();
    }

    @Override
    public Result<List<AiBotRecentlyUsed>> list(AiBotRecentlyUsed entity, Boolean asTree, String sortKey, String sortType) {
        LoginAccount account = SaTokenUtil.getLoginAccount();
        entity.setCreatedBy(account.getId());
        return super.list(entity, asTree, sortKey, sortType);
    }

    @Override
    protected Result<?> onSaveOrUpdateBefore(AiBotRecentlyUsed entity, boolean isSave) {
        entity.setCreated(new Date());
        entity.setCreatedBy(SaTokenUtil.getLoginAccount().getId());
        return super.onSaveOrUpdateBefore(entity, isSave);
    }
}