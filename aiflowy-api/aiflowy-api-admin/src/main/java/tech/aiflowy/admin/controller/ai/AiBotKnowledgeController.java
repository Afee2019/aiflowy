package tech.aiflowy.admin.controller.ai;

import com.mybatisflex.core.table.TableInfo;
import com.mybatisflex.core.table.TableInfoFactory;
import org.springframework.web.bind.annotation.PostMapping;
import tech.aiflowy.ai.entity.AiBotKnowledge;
import tech.aiflowy.ai.service.AiBotKnowledgeService;
import tech.aiflowy.common.annotation.UsePermission;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.entity.LoginAccount;
import tech.aiflowy.common.satoken.util.SaTokenUtil;
import tech.aiflowy.common.web.controller.BaseCurdController;
import com.mybatisflex.core.query.QueryWrapper;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import tech.aiflowy.common.web.jsonbody.JsonBody;

import java.math.BigInteger;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 控制层。
 *
 * @author michael
 * @since 2024-08-28
 */
@RestController
@RequestMapping("/api/v1/aiBotKnowledge")
@UsePermission(moduleName = "/api/v1/aiBot")
public class AiBotKnowledgeController extends BaseCurdController<AiBotKnowledgeService, AiBotKnowledge> {
    public AiBotKnowledgeController(AiBotKnowledgeService service) {
        super(service);
    }

    @GetMapping("list")
    @Override
    public Result<List<AiBotKnowledge>> list(AiBotKnowledge entity, Boolean asTree, String sortKey, String sortType) {
        QueryWrapper queryWrapper = QueryWrapper.create(entity, buildOperators(entity));
        queryWrapper.orderBy(buildOrderBy(sortKey, sortType, getDefaultOrderBy()));
        List<AiBotKnowledge> aiBotKnowledges = service.getMapper().selectListWithRelationsByQuery(queryWrapper);
        return Result.ok(aiBotKnowledges);
    }

    @PostMapping("updateBotKnowledgeIds")
    public Result<?> save(@JsonBody("botId") BigInteger botId, @JsonBody("knowledgeIds") BigInteger [] knowledgeIds) {
         service.saveBotAndKnowledge(botId, knowledgeIds);
        return Result.ok();
    }
}