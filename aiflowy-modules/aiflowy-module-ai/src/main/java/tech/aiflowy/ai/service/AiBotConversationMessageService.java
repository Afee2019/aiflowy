package tech.aiflowy.ai.service;

import com.mybatisflex.core.service.IService;
import dev.tinyflow.core.llm.Llm;
import tech.aiflowy.ai.entity.AiBotConversationMessage;
import tech.aiflowy.common.domain.Result;

import java.math.BigInteger;

/**
 *  服务层。
 *
 * @author Administrator
 * @since 2025-04-15
 */
public interface AiBotConversationMessageService extends IService<AiBotConversationMessage> {

    Result deleteConversation(String botId, String sessionId);

    Result updateConversation(String botId, String sessionId, String title);

    Result externalList(BigInteger botId);

    Boolean needRefreshConversationTitle(String sessionId, String userPrompt, Llm llm, BigInteger botId, int isExternalMsg);

}
