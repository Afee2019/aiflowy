package tech.aiflowy.common.cache;

import com.alicp.jetcache.Cache;
import com.alicp.jetcache.CacheManager;
import com.alicp.jetcache.anno.CacheType;
import com.alicp.jetcache.template.QuickConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.annotation.PostConstruct;

@Configuration
public class CacheConfig {

    @Autowired
    private CacheManager cacheManager;
    @Value("${jetcache.cacheType}")
    private String cacheType;

    private Cache<String, Object> defaultCache;

    @PostConstruct
    public void init() {
        CacheType type = CacheType.LOCAL;
        if ("remote".equals(cacheType)) {
            type = CacheType.REMOTE;
        }
        if ("both".equals(cacheType)) {
            type = CacheType.BOTH;
        }
        QuickConfig quickConfig = QuickConfig.newBuilder("")
                .cacheType(type)
                .build();
        defaultCache = cacheManager.getOrCreateCache(quickConfig);
    }

    @Bean("defaultCache")
    public Cache<String, Object> getDefaultCache() {
        return defaultCache;
    }
}
