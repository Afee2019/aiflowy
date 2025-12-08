package tech.aiflowy.ai.tinyflow.listener;

import dev.tinyflow.core.chain.Chain;
import dev.tinyflow.core.chain.Node;
import dev.tinyflow.core.chain.listener.NodeErrorListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class NodeErrorListenerForSave implements NodeErrorListener {

    private static final Logger log = LoggerFactory.getLogger(NodeErrorListenerForSave.class);

    @Override
    public void onError(Throwable error, Node node, Map<String, Object> nodeResult, Chain chain) {
        log.error("NodeErrorListenerForFront: {}", node, error);
    }
}
