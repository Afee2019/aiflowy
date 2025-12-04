package tech.aiflowy.ai.service.impl;

import tech.aiflowy.ai.entity.AiBotKnowledge;
import tech.aiflowy.ai.entity.AiBotPlugins;
import tech.aiflowy.ai.mapper.AiBotKnowledgeMapper;
import tech.aiflowy.ai.service.AiBotKnowledgeService;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;
import com.mybatisflex.core.query.QueryWrapper;

/**
 *  服务层实现。
 *
 * @author michael
 * @since 2024-08-28
 */
@Service
public class AiBotKnowledgeServiceImpl extends ServiceImpl<AiBotKnowledgeMapper, AiBotKnowledge> implements AiBotKnowledgeService {

    @Override
    public List<AiBotKnowledge> listByBotId(BigInteger botId) {

        QueryWrapper queryWrapper = QueryWrapper.create();
        queryWrapper.eq("bot_id",botId);

        return list(queryWrapper);
    }

    @Override
    public void saveBotAndKnowledge(BigInteger botId, BigInteger[] knowledgeIds) {
        this.remove(QueryWrapper.create().eq(AiBotKnowledge::getBotId, botId));
        List<AiBotKnowledge> list = new ArrayList<>(knowledgeIds.length);
        for (BigInteger knowledgeId : knowledgeIds) {
            AiBotKnowledge aiBotKnowledge = new AiBotKnowledge();
            aiBotKnowledge.setBotId(botId);
            aiBotKnowledge.setKnowledgeId(knowledgeId);
            list.add(aiBotKnowledge);
        }
        this.saveBatch(list);
    }
}
