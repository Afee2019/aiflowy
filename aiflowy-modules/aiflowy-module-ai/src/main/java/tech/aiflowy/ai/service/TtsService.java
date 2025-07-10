package tech.aiflowy.ai.service;

import java.util.concurrent.CompletableFuture;
import java.util.function.BiConsumer;

public interface TtsService {

    CompletableFuture<Void> streamTextToSpeech(String text, BiConsumer<String, Boolean> audioDataCallback,String sessionId);
}
