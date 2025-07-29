package tech.aiflowy.ai.service.impl;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.util.StringUtils;
import tech.aiflowy.ai.entity.AiBot;
import tech.aiflowy.ai.mapper.AiBotMapper;
import tech.aiflowy.ai.service.AiBotService;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;
import tech.aiflowy.common.domain.Result;
import tech.aiflowy.common.web.exceptions.BusinessException;
import tech.aiflowy.ai.utils.RegexUtils;
import com.mybatisflex.core.query.QueryWrapper;

import java.util.Map;

/**
 *  服务层实现。
 *
 * @author michael
 * @since 2024-08-23
 */
@Service
public class AiBotServiceImpl extends ServiceImpl<AiBotMapper, AiBot> implements AiBotService {

    private static final Logger log = LoggerFactory.getLogger(AiBotServiceImpl.class);

    @Override
    public Result getDetail(String id) {
        AiBot aiBot = null;

        if (id.matches(RegexUtils.ALL_NUMBER)){
            aiBot = getById(id);

            if (aiBot == null) {
               aiBot = getByAlias(id);
            }

        }else{
            aiBot = getByAlias(id);
        }

        return Result.success(aiBot);
    }

    @Override
    public AiBot getByAlias(String alias){
        QueryWrapper queryWrapper = QueryWrapper.create();
        queryWrapper.eq("alias",alias);
        return getOne(queryWrapper);
    }

    @Override
    public void updateBotLlmId(AiBot aiBot) {
        AiBot byId = getById(aiBot.getId());

        if (byId == null) {
            log.error("修改bot的llmId失败，bot不存在！");
            throw new BusinessException("bot不存在！");
        }

        byId.setLlmId(aiBot.getLlmId());

        updateById(byId,false);

    }
}
