package tech.aiflowy.ai.tinyflow.config;

import dev.tinyflow.core.parser.ChainParser;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import tech.aiflowy.ai.tinyflow.service.TinyFlowConfigService;

import javax.annotation.Resource;

@Configuration
public class ChainParserConfig {

    @Resource
    private TinyFlowConfigService tinyFlowConfigService;

    @Bean
    public ChainParser chainParser() {
        ChainParser chainParser = ChainParser.builder()
                .withDefaultParsers(true)
                .build();
        tinyFlowConfigService.initProvidersAndNodeParsers(chainParser);
        return chainParser;
    }
}
