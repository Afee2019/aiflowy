package tech.aiflowy.ai.service;

import com.mybatisflex.core.service.IService;
import tech.aiflowy.ai.entity.AiBotPlugins;
import tech.aiflowy.ai.entity.AiPlugin;
import tech.aiflowy.common.domain.Result;

import java.math.BigInteger;
import java.util.List;

/**
 *  服务层。
 *
 * @author michael
 * @since 2025-04-07
 */
public interface AiBotPluginsService extends IService<AiBotPlugins> {

    List<AiPlugin> getList(String botId);

    boolean doRemove(String botId, String pluginId);

    List<BigInteger> getBotPluginToolIds(String botId);

    void saveBotAndPluginTool(BigInteger botId, BigInteger[] pluginToolIds);
}
