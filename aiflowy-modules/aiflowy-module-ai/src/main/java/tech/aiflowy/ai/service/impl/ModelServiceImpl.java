
package tech.aiflowy.ai.service.impl;

import cn.hutool.core.util.StrUtil;
import com.agentsflex.core.document.Document;
import com.agentsflex.core.model.chat.ChatModel;
import com.agentsflex.core.model.embedding.EmbeddingModel;
import com.agentsflex.core.model.rerank.RerankModel;
import com.agentsflex.core.store.VectorData;
import com.alicp.jetcache.Cache;
import com.mybatisflex.core.query.QueryWrapper;
import com.mybatisflex.spring.service.impl.ServiceImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.CollectionUtils;
import tech.aiflowy.ai.entity.Model;
import tech.aiflowy.ai.entity.ModelProvider;
import tech.aiflowy.ai.mapper.ModelMapper;
import tech.aiflowy.ai.service.ModelProviderService;
import tech.aiflowy.ai.service.ModelService;
import tech.aiflowy.common.web.exceptions.BusinessException;

import javax.annotation.Resource;
import java.math.BigInteger;
import java.util.*;
import java.util.stream.Collectors;

/**
 * 服务层实现。
 *
 * @author michael
 * @since 2024-08-23
 */
@Service
public class ModelServiceImpl extends ServiceImpl<ModelMapper, Model> implements ModelService {

    @Autowired
    ModelMapper modelMapper;

    @Autowired
    ModelProviderService llmProviderService;

    @Resource
    private Cache<String, Object> cache;


    @Override
    public boolean addAiLlm(Model entity) {
        int insert = modelMapper.insert(entity);
        if (insert <= 0) {
            return false;
        }
        return true;
    }

    private static final Logger log = LoggerFactory.getLogger(ModelServiceImpl.class);

    @Override
    public void verifyLlmConfig(Model llm) {
        String modelType = llm.getModelType();
        if ("chatModel".equals(modelType)) {
            // 走聊天验证逻辑
            verifyChatLlm(llm);
            return;
        }

        if ("embeddingModel".equals(modelType)) {
            // 走向量化验证逻辑
            verifyEmbedLlm(llm);
            return;

        }

        if ("rerankModel".equals(modelType)) {
            // 走重排验证逻辑
            verifyRerankLlm(llm);
            return;

        }

        // 以上不满足，视为验证失败
        throw new BusinessException("校验失败！");

    }

    @Override
    public Map<String, Map<String, List<Model>>> getList(Model entity) {
        String[] llmModelTypes = {"chatModel", "embeddingModel", "rerankModel"};
        Map<String, Map<String, List<Model>>> result = new HashMap<>();

        QueryWrapper queryWrapper = new QueryWrapper()
                .eq(Model::getProviderId, entity.getProviderId());
        queryWrapper.eq(Model::getWithUsed, entity.getWithUsed());
        List<Model> totalList = modelMapper.selectListWithRelationsByQuery(queryWrapper);
        for (String modelType : llmModelTypes) {
            Map<String, List<Model>> groupMap = groupLlmByGroupName(totalList, modelType);
            if (!CollectionUtils.isEmpty(groupMap)) {
                result.put(modelType, groupMap);
            }
        }

        return result;
    }

    private Map<String, List<Model>> groupLlmByGroupName(List<Model> totalList, String targetModelType) {
        if (CollectionUtils.isEmpty(totalList)) {
            return Collections.emptyMap();
        }

        return totalList.stream()
                .filter(aiLlm -> targetModelType.equals(aiLlm.getModelType())
                        && aiLlm.getGroupName() != null)
                .collect(Collectors.groupingBy(Model::getGroupName));
    }


    private void verifyRerankLlm(Model llm) {
        RerankModel rerankModel = llm.toRerankModel();
        List<Document> documents = new ArrayList<>();
        documents.add(Document.of("Paris is the capital of France."));
        documents.add(Document.of("London is the capital of England."));
        documents.add(Document.of("Tokyo is the capital of Japan."));
        documents.add(Document.of("Beijing is the capital of China."));
        documents.add(Document.of("Washington, D.C. is the capital of the United States."));
        documents.add(Document.of("Moscow is the capital of Russia."));
        try {
            List<Document> rerank = rerankModel.rerank("What is the capital of France?", documents);
            if (rerank == null || rerank.isEmpty()) {
                throw new BusinessException("校验未通过，请前往后端日志查看详情！");
            }
        } catch (Exception e) {
            log.error("校验失败：{}", e.getMessage());
            throw new BusinessException("校验未通过，请前往后端日志查看详情！");
        }
    }

    private void verifyEmbedLlm(Model llm) {
        try {
            EmbeddingModel transLlm = llm.toEmbeddingModel();
            VectorData vectorData = transLlm.embed("这是一条校验模型配置的文本");
            if (vectorData.getVector() == null) {
                throw new BusinessException("校验未通过，请前往后端日志查看详情！");
            }
            log.info("取到向量数据，校验结果通过");
        } catch (Exception e) {
            log.error("模型配置校验失败:{}", e.getMessage());
            throw new BusinessException("校验未通过，请前往后端日志查看详情！");
        }
    }

    private void verifyChatLlm(Model llm) {

        ChatModel chatModel = llm.toChatModel();
        if (chatModel == null) {
            throw new BusinessException("chatModel为空");
        }
        try {
            String response = chatModel.chat("我在对模型配置进行校验，你收到这条消息无需做任何思考，直接回复一个“你好”即可!");
            if (response == null) {
                throw new BusinessException("校验未通过，请前往后端日志查看详情！");
            }
            log.info("校验结果：{}", response);
        } catch (Exception e) {
            log.error("校验失败：{}", e.getMessage());
            throw new BusinessException("校验未通过，请前往后端日志查看详情！");
        }


        Map<String, Object> options = llm.getOptions();
//        if (options != null && options.get("multimodal") != null && (boolean) options.get("multimodal")) {
//
//            textPrompt = new ImagePrompt("我在对模型配置进行校验，你无需描述图片，只需回答“看到了图片”即可",
//                "http://115.190.9.61:7900/aiflowy-pro/public/aibot/files/40b64e32b081942bd7ab30f8a369f2a34fc7fafc04f45c50cd96d8a102fd7afa.jpg");
//
//        } else {
//            textPrompt = new TextPrompt("我在对模型配置进行校验，你收到这条消息无需做任何思考，直接回复一个“你好”即可!");
//        }


    }

    @Override
    public void removeByEntity(Model entity) {
        QueryWrapper queryWrapper = QueryWrapper.create().eq(Model::getProviderId, entity.getProviderId()).eq(Model::getGroupName, entity.getGroupName());
        modelMapper.deleteByQuery(queryWrapper);
    }

    @Override
    public Model getLlmInstance(BigInteger llmId) {
        Model aillm = getById(llmId);
        if (aillm == null) {
            return null;
        }
        ModelProvider provider = llmProviderService.getById(aillm.getProviderId());
        if (provider == null) {
            return aillm;
        }
        aillm.setModelProvider(provider);
        if (StrUtil.isBlank(aillm.getApiKey())) {
            aillm.setApiKey(provider.getApiKey());
        }
        if (StrUtil.isBlank(aillm.getEndpoint())) {
            aillm.setEndpoint(provider.getEndpoint());
        }

        Map<String, Object> options = aillm.getOptions();
        if (options == null) {
            options = new HashMap<>();
            aillm.setOptions(options);
        }

        String chatPath = (String) options.get("chatPath");
        if (StrUtil.isBlank(chatPath)) {
            options.put("chatPath", provider.getChatPath());
        }

        String embedPath = (String) options.get("embedPath");
        if (StrUtil.isBlank(embedPath)) {
            options.put("embedPath", provider.getEmbedPath());
        }

        String rerankPath = (String) options.get("rerankPath");
        if (StrUtil.isBlank(rerankPath)) {
            options.put("rerankPath", provider.getRerankPath());
        }

        String llmEndpoint = (String) options.get("llmEndpoint");
        if (StrUtil.isBlank(llmEndpoint)) {
            options.put("llmEndpoint", provider.getEndpoint());
        }

        return aillm;
    }

    @Override
    public void updateByEntity(Model entity) {
        QueryWrapper queryWrapper = QueryWrapper.create().eq(Model::getProviderId, entity.getProviderId())
                .eq(Model::getGroupName, entity.getGroupName());
        Model model = new Model();
        model.setWithUsed(entity.getWithUsed());
        modelMapper.updateByQuery(model, queryWrapper);
    }
}
