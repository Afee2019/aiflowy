package tech.aiflowy.usercenter.controller.ai;

import cn.dev33.satoken.annotation.SaCheckPermission;
import cn.dev33.satoken.stp.StpUtil;
import dev.tinyflow.core.chain.*;
import dev.tinyflow.core.chain.repository.ChainStateRepository;
import dev.tinyflow.core.chain.repository.NodeStateRepository;
import dev.tinyflow.core.chain.runtime.ChainExecutor;
import dev.tinyflow.core.parser.ChainParser;
import org.springframework.web.bind.annotation.*;
import tech.aiflowy.ai.entity.AiWorkflow;
import tech.aiflowy.ai.service.AiWorkflowService;
import tech.aiflowy.ai.tinyflow.entity.ChainInfo;
import tech.aiflowy.ai.tinyflow.entity.NodeInfo;
import tech.aiflowy.ai.tinyflow.service.TinyFlowService;
import tech.aiflowy.common.annotation.UsePermission;
import tech.aiflowy.common.constant.Constants;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.satoken.util.SaTokenUtil;
import tech.aiflowy.common.web.controller.BaseCurdController;
import tech.aiflowy.common.web.jsonbody.JsonBody;

import javax.annotation.Resource;
import java.io.Serializable;
import java.math.BigInteger;
import java.util.*;

/**
 * 工作流
 */
@RestController
@RequestMapping("/userCenter/aiWorkflow")
@UsePermission(moduleName = "/api/v1/aiWorkflow")
public class UcAiWorkflowController extends BaseCurdController<AiWorkflowService, AiWorkflow> {

    @Resource
    private ChainExecutor chainExecutor;
    @Resource
    private ChainParser chainParser;
    @Resource
    private TinyFlowService tinyFlowService;

    public UcAiWorkflowController(AiWorkflowService service) {
        super(service);
    }

    /**
     * 节点单独运行
     */
    @PostMapping("/singleRun")
    @SaCheckPermission("/api/v1/aiWorkflow/save")
    public Result<?> singleRun(
            @JsonBody(value = "workflowId", required = true) BigInteger workflowId,
            @JsonBody(value = "nodeId", required = true) String nodeId,
            @JsonBody("variables") Map<String, Object> variables) {

        AiWorkflow workflow = service.getById(workflowId);
        if (workflow == null) {
            return Result.fail(1, "工作流不存在");
        }
        Map<String, Object> res = chainExecutor.executeNode(workflowId.toString(), nodeId, variables);
        return Result.ok(res);
    }

    /**
     * 运行工作流 - v2
     */
    @PostMapping("/runAsync")
    @SaCheckPermission("/api/v1/aiWorkflow/save")
    public Result<String> runAsync(@JsonBody(value = "id", required = true) BigInteger id,
                                   @JsonBody("variables") Map<String, Object> variables) {
        if (variables == null) {
            variables = new HashMap<>();
        }
        AiWorkflow workflow = service.getById(id);
        if (workflow == null) {
            throw new RuntimeException("工作流不存在");
        }
        if (StpUtil.isLogin()) {
            variables.put(Constants.LOGIN_USER_KEY, SaTokenUtil.getLoginAccount());
        }
        String executeId = chainExecutor.executeAsync(id.toString(), variables);
        return Result.ok(executeId);
    }

    /**
     * 获取工作流运行状态 - v2
     */
    @PostMapping("/getChainStatus")
    public Result<ChainInfo> getChainStatus(@JsonBody(value = "executeId") String executeId,
                                            @JsonBody("nodes") List<NodeInfo> nodes) {
        ChainInfo res = tinyFlowService.getChainStatus(executeId, nodes);
        return Result.ok(res);
    }

    /**
     * 恢复工作流运行 - v2
     */
    @PostMapping("/resume")
    @SaCheckPermission("/api/v1/aiWorkflow/save")
    public Result<Void> resume(@JsonBody(value = "executeId", required = true) String executeId,
                               @JsonBody("confirmParams") Map<String, Object> confirmParams) {
        chainExecutor.resumeAsync(executeId, confirmParams);
        return Result.ok();
    }

    /**
     * 获取工作流参数 - v2
     */
    @GetMapping("getRunningParameters")
    @SaCheckPermission("/api/v1/aiWorkflow/query")
    public Result<?> getRunningParameters(@RequestParam BigInteger id) {
        AiWorkflow workflow = service.getById(id);

        if (workflow == null) {
            return Result.fail(1, "can not find the workflow by id: " + id);
        }

        ChainDefinition definition = chainParser.parse(workflow.getContent());
        if (definition == null) {
            return Result.fail(2, "节点配置错误，请检查! ");
        }
        List<Parameter> chainParameters = definition.getStartParameters();
        Map<String, Object> res = new HashMap<>();
        res.put("parameters", chainParameters);
        res.put("title", workflow.getTitle());
        res.put("description", workflow.getDescription());
        res.put("icon", workflow.getIcon());
        return Result.ok(res);
    }

    @Override
    protected Result<?> onSaveOrUpdateBefore(AiWorkflow entity, boolean isSave) {
        return Result.fail("-");
    }

    @Override
    protected Result<?> onRemoveBefore(Collection<Serializable> ids) {
        return Result.fail("-");
    }
}
