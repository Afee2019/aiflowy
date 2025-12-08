package tech.aiflowy.ai.tinyflow.listener;

import dev.tinyflow.core.chain.Chain;
import dev.tinyflow.core.chain.Event;
import dev.tinyflow.core.chain.event.*;
import dev.tinyflow.core.chain.listener.ChainEventListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ChainEventListenerForSave implements ChainEventListener {


    private static final Logger log = LoggerFactory.getLogger(ChainEventListenerForSave.class);

    @Override
    public void onEvent(Event event, Chain chain) {
        if (event instanceof ChainStartEvent) {
            handleChainStartEvent((ChainStartEvent) event, chain);
        }
        if (event instanceof ChainEndEvent) {
            handleChainEndEvent((ChainEndEvent) event, chain);
        }
        if (event instanceof NodeStartEvent) {
            handleNodeStartEvent((NodeStartEvent) event, chain);
        }
        if (event instanceof NodeEndEvent) {
            handleNodeEndEvent((NodeEndEvent) event, chain);
        }
        if (event instanceof ChainStatusChangeEvent) {
            handleChainStatusChangeEvent((ChainStatusChangeEvent) event, chain);
        }
        if (event instanceof ChainResumeEvent) {
            handleChainResumeEvent((ChainResumeEvent) event, chain);
        }
    }

    private void handleChainStartEvent(ChainStartEvent event, Chain chain) {
        log.info("ChainStartEvent: {}", event);
    }

    private void handleChainEndEvent(ChainEndEvent event, Chain chain) {
        log.info("ChainEndEvent: {}", event);
    }

    private void handleNodeStartEvent(NodeStartEvent event, Chain chain) {
        log.info("NodeStartEvent: {}", event);
    }

    private void handleNodeEndEvent(NodeEndEvent event, Chain chain) {
        log.info("NodeEndEvent: {}", event);
    }

    private void handleChainStatusChangeEvent(ChainStatusChangeEvent event, Chain chain) {
        log.info("ChainStatusChangeEvent: {}", event);
    }

    private void handleChainResumeEvent(ChainResumeEvent event, Chain chain) {
        log.info("ChainResumeEvent: {}", event);
    }

}
