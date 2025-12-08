package tech.aiflowy.ai.tinyflow.listener;

import dev.tinyflow.core.chain.Chain;
import dev.tinyflow.core.chain.listener.ChainErrorListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ChainErrorListenerForSave implements ChainErrorListener {

    private static final Logger log = LoggerFactory.getLogger(ChainErrorListenerForSave.class);

    @Override
    public void onError(Throwable error, Chain chain) {
        log.error("ChainErrorListenerForFront: {}", chain, error);
    }
}
