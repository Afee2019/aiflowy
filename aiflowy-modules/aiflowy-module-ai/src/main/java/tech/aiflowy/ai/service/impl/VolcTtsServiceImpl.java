package tech.aiflowy.ai.service.impl;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import okhttp3.*;
//import okhttp3.logging.HttpLoggingInterceptor;
import okio.ByteString;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import tech.aiflowy.ai.service.TtsService;

import java.io.IOException;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.function.BiConsumer;

@Service("VolcTtsService")
public class VolcTtsServiceImpl implements TtsService {

    private static final Logger logger = LoggerFactory.getLogger(VolcTtsServiceImpl.class);
    private final ObjectMapper objectMapper = new ObjectMapper();

    // 协议常量
    private static final int PROTOCOL_VERSION = 0b0001;
    private static final int DEFAULT_HEADER_SIZE = 0b0001;

    // Message Type:
    private static final int FULL_CLIENT_REQUEST = 0b0001;
    private static final int AUDIO_ONLY_RESPONSE = 0b1011;
    private static final int FULL_SERVER_RESPONSE = 0b1001;
    private static final int ERROR_INFORMATION = 0b1111;

    // Message Type Specific Flags
    private static final int MsgTypeFlagNoSeq = 0b0000;
    private static final int MsgTypeFlagPositiveSeq = 0b1;
    private static final int MsgTypeFlagLastNoSeq = 0b10;
    private static final int MsgTypeFlagNegativeSeq = 0b11;
    private static final int MsgTypeFlagWithEvent = 0b100;

    // Message Serialization
    private static final int NO_SERIALIZATION = 0b0000;
    private static final int JSON = 0b0001;

    // Message Compression
    private static final int COMPRESSION_NO = 0b0000;
    private static final int COMPRESSION_GZIP = 0b0001;

    // 事件常量
    public static final int EVENT_NONE = 0;
    public static final int EVENT_Start_Connection = 1;
    public static final int EVENT_FinishConnection = 2;
    public static final int EVENT_ConnectionStarted = 50;
    public static final int EVENT_ConnectionFailed = 51;
    public static final int EVENT_ConnectionFinished = 52;
    public static final int EVENT_StartSession = 100;
    public static final int EVENT_FinishSession = 102;
    public static final int EVENT_SessionStarted = 150;
    public static final int EVENT_SessionFinished = 152;
    public static final int EVENT_SessionFailed = 153;
    public static final int EVENT_TaskRequest = 200;
    public static final int EVENT_TTSSentenceStart = 350;
    public static final int EVENT_TTSSentenceEnd = 351;
    public static final int EVENT_TTSResponse = 352;

    // TTS配置
    @Value("${voiceInput.volcengine.app.appId}")
    private String appId ; // 从配置获取
    @Value("${voiceInput.volcengine.app.token}")
    private String token ; // 从配置获取
    private String speaker = "zh_female_shuangkuaisisi_moon_bigtts";
    private String url = "wss://openspeech.bytedance.com/api/v3/tts/bidirection";

    /**
     * 流式文本转语音
     * @param text 要转换的文本
     * @param audioDataCallback 音频数据回调 (base64Data, isComplete)
     * @return CompletableFuture<Void>
     */
    public CompletableFuture<Void> streamTextToSpeech(String text, BiConsumer<String, Boolean> audioDataCallback,String connectId) {
        CompletableFuture<Void> future = new CompletableFuture<>();

        if ("${{over}}$".equalsIgnoreCase(text)) {
            audioDataCallback.accept("", true);
            future.complete(null);
            return future;
        }

        final Request request = new Request.Builder()
                .url(url)
                .header("X-Api-App-Key", appId)
                .header("X-Api-Access-Key", token)
                .header("X-Api-Resource-Id", "volc.service_type.10029")
                .header("X-Api-Connect-Id", connectId)
                .build();

//        HttpLoggingInterceptor loggingInterceptor = new HttpLoggingInterceptor();
//        loggingInterceptor.setLevel(HttpLoggingInterceptor.Level.HEADERS);

        final OkHttpClient okHttpClient = new OkHttpClient.Builder()
                .pingInterval(50, TimeUnit.SECONDS)
                .readTimeout(300, TimeUnit.SECONDS)
                .writeTimeout(300, TimeUnit.SECONDS)
                .build();

        okHttpClient.newWebSocket(request, new WebSocketListener() {
            final String sessionId = UUID.randomUUID().toString().replace("-", "");

            @Override
            public void onOpen(WebSocket webSocket, Response response) {
                logger.info("TTS WebSocket连接建立，Logid: {}", response.header("X-Tt-Logid"));
                startConnection(webSocket);
            }

            @Override
            public void onMessage(WebSocket webSocket, ByteString bytes) {
                TTSResponse response = parserResponse(bytes.toByteArray());
                if (response == null) return;

                logger.debug("TTS响应: event={}", response.optional.event);

                switch (response.optional.event) {
                    case EVENT_ConnectionFailed:
                    case EVENT_SessionFailed:
                        logger.error("TTS连接或会话失败");
                        future.completeExceptionally(new RuntimeException("TTS连接失败"));
                        return;

                    case EVENT_ConnectionStarted:
                        startTTSSession(webSocket, sessionId, speaker);
                        break;

                    case EVENT_SessionStarted:
                        sendMessage(webSocket, speaker, sessionId, text);
                        finishSession(webSocket, sessionId);
                        break;

                    case EVENT_TTSSentenceStart:
                    case EVENT_TTSSentenceEnd:
                        break;

                    case EVENT_TTSResponse:
                        if (response.payload != null && response.header.message_type == AUDIO_ONLY_RESPONSE) {
                            // 将音频数据转为base64并回调
                            String base64Audio = java.util.Base64.getEncoder().encodeToString(response.payload);
                            audioDataCallback.accept(base64Audio, false);
                        }
                        break;

                    case EVENT_SessionFinished:
                        finishConnection(webSocket);
                        break;

                    case EVENT_ConnectionFinished:
                        logger.info("TTS连接结束，connectId: {}", connectId);
                        audioDataCallback.accept("", true); // 标记完成
                        future.complete(null);
                        webSocket.close(1000,"正常关闭");
                        break;

                    default:
                        break;
                }
            }

            @Override
            public void onFailure(WebSocket webSocket, Throwable t, Response response) {
                logger.error("TTS WebSocket连接失败", t);
                future.completeExceptionally(t);
            }
        });

        return future;
    }


    private int bytesToInt(byte[] src) {
        if (src == null || (src.length != 4)) {
            throw new IllegalArgumentException("");
        }
        return ((src[0] & 0xFF) << 24)
                | ((src[1] & 0xff) << 16)
                | ((src[2] & 0xff) << 8)
                | ((src[3] & 0xff));
    }

    private byte[] intToBytes(int a) {
        return new byte[]{
                (byte) ((a >> 24) & 0xFF),
                (byte) ((a >> 16) & 0xFF),
                (byte) ((a >> 8) & 0xFF),
                (byte) (a & 0xFF)
        };
    }

    public static class Header {
        public int protocol_version = PROTOCOL_VERSION;
        public int header_size = DEFAULT_HEADER_SIZE;
        public int message_type;
        public int message_type_specific_flags = MsgTypeFlagWithEvent;
        public int serialization_method = NO_SERIALIZATION;
        public int message_compression = COMPRESSION_NO;
        public int reserved = 0;

        public Header() {}

        public Header(int protocol_version, int header_size, int message_type, int message_type_specific_flags,
                      int serialization_method, int message_compression, int reserved) {
            this.protocol_version = protocol_version;
            this.header_size = header_size;
            this.message_type = message_type;
            this.message_type_specific_flags = message_type_specific_flags;
            this.serialization_method = serialization_method;
            this.message_compression = message_compression;
            this.reserved = reserved;
        }

        public byte[] getBytes() {
            return new byte[]{
                    (byte) ((protocol_version << 4) | header_size),
                    (byte) (message_type << 4 | message_type_specific_flags),
                    (byte) ((serialization_method << 4) | message_compression),
                    (byte) reserved
            };
        }
    }

    public static class Optional {
        public int event = EVENT_NONE;
        public String sessionId;
        public int errorCode;
        public String connectionId;
        public String response_meta_json;

        public Optional(int event, String sessionId) {
            this.event = event;
            this.sessionId = sessionId;
        }

        public Optional() {}

        public byte[] getBytes() {
            byte[] bytes = new byte[0];
            if (event != EVENT_NONE) {
                bytes = intToBytes(event);
            }
            if (sessionId != null) {
                byte[] sessionIdSize = intToBytes(sessionId.getBytes().length);
                final byte[] temp = bytes;
                int desPos = 0;
                bytes = new byte[temp.length + sessionIdSize.length + sessionId.getBytes().length];
                System.arraycopy(temp, 0, bytes, desPos, temp.length);
                desPos += temp.length;
                System.arraycopy(sessionIdSize, 0, bytes, desPos, sessionIdSize.length);
                desPos += sessionIdSize.length;
                System.arraycopy(sessionId.getBytes(), 0, bytes, desPos, sessionId.getBytes().length);
            }
            return bytes;
        }

        private byte[] intToBytes(int a) {
            return new byte[]{
                    (byte) ((a >> 24) & 0xFF),
                    (byte) ((a >> 16) & 0xFF),
                    (byte) ((a >> 8) & 0xFF),
                    (byte) (a & 0xFF)
            };
        }
    }

    private static class Pair<F, S> {
        public Pair(F fst, S snd) {
            this.fst = fst;
            this.snd = snd;
        }
        public F fst;
        public S snd;
    }

    public static class TTSResponse {
        public Header header;
        public Optional optional;
        public int payloadSize;
        transient public byte[] payload;
        public String payloadJson;

        @Override
        public String toString() {
            return "TTSResponse{event=" + (optional != null ? optional.event : "null") + ", payloadSize=" + payloadSize + "}";
        }
    }

    // 完全保留所有解析方法，改为实例方法
    private TTSResponse parserResponse(byte[] res) {
        if (res == null || res.length == 0) {
            return null;
        }
        final TTSResponse response = new TTSResponse();
        Header header = new Header();
        response.header = header;

        final byte num = 0b00001111;
        header.protocol_version = (res[0] >> 4) & num;
        header.header_size = res[0] & 0x0f;
        header.message_type = (res[1] >> 4) & num;
        header.message_type_specific_flags = res[1] & 0x0f;
        header.serialization_method = res[2] >> num;
        header.message_compression = res[2] & 0x0f;
        header.reserved = res[3];

        int offset = 4;
        response.optional = new Optional();

        if (header.message_type == FULL_SERVER_RESPONSE || header.message_type == AUDIO_ONLY_RESPONSE) {
            offset = readEvent(res, header.message_type_specific_flags, response, offset);
            final int event = response.optional.event;

            switch (event) {
                case EVENT_ConnectionStarted:
                    readConnectStarted(res, response, offset);
                    break;
                case EVENT_ConnectionFailed:
                    readConnectFailed(res, response, offset);
                    break;
                case EVENT_SessionStarted:
                    offset = readSessionId(res, response, offset);
                    break;
                case EVENT_TTSResponse:
                    offset = readSessionId(res, response, offset);
                    offset = readPayload(res, response, offset);
                    break;
                case EVENT_TTSSentenceStart:
                case EVENT_TTSSentenceEnd:
                    offset = readSessionId(res, response, offset);
                    offset = readPayloadJson(res, response, offset);
                    break;
                case EVENT_SessionFailed:
                case EVENT_SessionFinished:
                    offset = readSessionId(res, response, offset);
                    readMetaJson(res, response, offset);
                    break;
                default:
                    break;
            }
        } else if (header.message_type == ERROR_INFORMATION) {
            offset = readErrorCode(res, response, offset);
            readPayload(res, response, offset);
        }
        return response;
    }

    private void readConnectStarted(byte[] res, TTSResponse response, int start) {
        start = readConnectId(res, response, start);
    }

    private void readConnectFailed(byte[] res, TTSResponse response, int start) {
        start = readConnectId(res, response, start);
        readMetaJson(res, response, start);
    }

    private int readConnectId(byte[] res, TTSResponse response, int start) {
        Pair<Integer, String> r = readText(res, start);
        start = r.fst;
        response.optional.connectionId = r.snd;
        return start;
    }

    private int readMetaJson(byte[] res, TTSResponse response, int start) {
        Pair<Integer, String> r = readText(res, start);
        start = r.fst;
        response.optional.response_meta_json = r.snd;
        return start;
    }

    private int readPayloadJson(byte[] res, TTSResponse response, int start) {
        Pair<Integer, String> r = readText(res, start);
        start = r.fst;
        response.payloadJson = r.snd;
        return start;
    }

    private Pair<Integer, String> readText(byte[] res, int start) {
        byte[] b = new byte[4];
        System.arraycopy(res, start, b, 0, b.length);
        start += b.length;
        int size = bytesToInt(b);
        b = new byte[size];
        System.arraycopy(res, start, b, 0, b.length);
        String text = new String(b);
        start += b.length;
        return new Pair<>(start, text);
    }

    private int readPayload(byte[] res, TTSResponse response, int start) {
        byte[] b = new byte[4];
        System.arraycopy(res, start, b, 0, b.length);
        start += b.length;
        int size = bytesToInt(b);
        response.payloadSize += size;
        b = new byte[size];
        System.arraycopy(res, start, b, 0, b.length);
        response.payload = b;
        start += b.length;
        return start;
    }

    private int readErrorCode(byte[] res, TTSResponse response, int start) {
        byte[] b = new byte[4];
        System.arraycopy(res, start, b, 0, b.length);
        start += b.length;
        response.optional.errorCode = bytesToInt(b);
        return start;
    }

    private int readEvent(byte[] res, int masTypeFlag, TTSResponse response, int start) {
        if (masTypeFlag == MsgTypeFlagWithEvent) {
            byte[] temp = new byte[4];
            System.arraycopy(res, start, temp, 0, temp.length);
            response.optional.event = bytesToInt(temp);
            start += temp.length;
        }
        return start;
    }

    private int readSessionId(byte[] res, TTSResponse response, int start) {
        Pair<Integer, String> r = readText(res, start);
        start = r.fst;
        response.optional.sessionId = r.snd;
        return start;
    }

    // 保留所有发送方法，改为实例方法
    private boolean startConnection(WebSocket webSocket) {
        byte[] header = new Header(
                PROTOCOL_VERSION,
                FULL_CLIENT_REQUEST,
                DEFAULT_HEADER_SIZE,
                MsgTypeFlagWithEvent,
                JSON,
                COMPRESSION_NO,
                0).getBytes();
        byte[] optional = new Optional(EVENT_Start_Connection, null).getBytes();
        byte[] payload = "{}".getBytes();
        return sendEvent(webSocket, header, optional, payload);
    }

    private boolean finishConnection(WebSocket webSocket) {
        byte[] header = new Header(
                PROTOCOL_VERSION,
                FULL_CLIENT_REQUEST,
                DEFAULT_HEADER_SIZE,
                MsgTypeFlagWithEvent,
                JSON,
                COMPRESSION_NO,
                0).getBytes();
        byte[] optional = new Optional(EVENT_FinishConnection, null).getBytes();
        byte[] payload = "{}".getBytes();
        return sendEvent(webSocket, header, optional, payload);
    }

    private boolean finishSession(WebSocket webSocket, String sessionId) {
        byte[] header = new Header(
                PROTOCOL_VERSION,
                FULL_CLIENT_REQUEST,
                DEFAULT_HEADER_SIZE,
                MsgTypeFlagWithEvent,
                JSON,
                COMPRESSION_NO,
                0).getBytes();
        byte[] optional = new Optional(EVENT_FinishSession, sessionId).getBytes();
        byte[] payload = "{}".getBytes();
        return sendEvent(webSocket, header, optional, payload);
    }

    private boolean startTTSSession(WebSocket webSocket, String sessionId, String speaker) {
        byte[] header = new Header(
                PROTOCOL_VERSION,
                FULL_CLIENT_REQUEST,
                DEFAULT_HEADER_SIZE,
                MsgTypeFlagWithEvent,
                JSON,
                COMPRESSION_NO,
                0).getBytes();

        final int event = EVENT_StartSession;
        byte[] optional = new Optional(event, sessionId).getBytes();

        try {
            ObjectNode payloadJObj = objectMapper.createObjectNode();
            ObjectNode user = objectMapper.createObjectNode();
            user.put("uid", "123456");

            payloadJObj.set("user", user);
            payloadJObj.put("event", event);
            payloadJObj.put("namespace", "BidirectionalTTS");

            ObjectNode req_params = objectMapper.createObjectNode();
            req_params.put("speaker", speaker);

            ObjectNode audio_params = objectMapper.createObjectNode();
            audio_params.put("format", "mp3");
            audio_params.put("sample_rate", 24000);
            audio_params.put("enable_timestamp", true);

            req_params.set("audio_params", audio_params);
            payloadJObj.set("req_params", req_params);

            byte[] payload = payloadJObj.toString().getBytes();
            return sendEvent(webSocket, header, optional, payload);
        } catch (Exception e) {
            logger.error("构建TTS Session请求失败", e);
            return false;
        }
    }

    private boolean sendMessage(WebSocket webSocket, String speaker, String sessionId, String text) {
        byte[] header = new Header(
                PROTOCOL_VERSION,
                FULL_CLIENT_REQUEST,
                DEFAULT_HEADER_SIZE,
                MsgTypeFlagWithEvent,
                JSON,
                COMPRESSION_NO,
                0).getBytes();

        byte[] optional = new Optional(EVENT_TaskRequest, sessionId).getBytes();

        try {
            ObjectNode payloadJObj = objectMapper.createObjectNode();
            ObjectNode user = objectMapper.createObjectNode();
            user.put("uid", "123456");
            payloadJObj.set("user", user);

            payloadJObj.put("event", EVENT_TaskRequest);
            payloadJObj.put("namespace", "BidirectionalTTS");

            ObjectNode req_params = objectMapper.createObjectNode();
            req_params.put("text", text);
            req_params.put("speaker", speaker);

            ObjectNode audio_params = objectMapper.createObjectNode();
            audio_params.put("format", "mp3");
            audio_params.put("sample_rate", 24000);

            req_params.set("audio_params", audio_params);
            payloadJObj.set("req_params", req_params);

            byte[] payload = payloadJObj.toString().getBytes();
            return sendEvent(webSocket, header, optional, payload);
        } catch (Exception e) {
            logger.error("构建TTS消息请求失败", e);
            return false;
        }
    }

    private boolean sendEvent(WebSocket webSocket, byte[] header, byte[] optional, byte[] payload) {
        assert webSocket != null;
        assert header != null;
        assert payload != null;

        final byte[] payloadSizeBytes = intToBytes(payload.length);
        byte[] requestBytes = new byte[
                header.length
                        + (optional == null ? 0 : optional.length)
                        + payloadSizeBytes.length + payload.length];
        int desPos = 0;
        System.arraycopy(header, 0, requestBytes, desPos, header.length);
        desPos += header.length;
        if (optional != null) {
            System.arraycopy(optional, 0, requestBytes, desPos, optional.length);
            desPos += optional.length;
        }
        System.arraycopy(payloadSizeBytes, 0, requestBytes, desPos, payloadSizeBytes.length);
        desPos += payloadSizeBytes.length;
        System.arraycopy(payload, 0, requestBytes, desPos, payload.length);
        return webSocket.send(ByteString.of(requestBytes));
    }

    // 配置方法
    public void setAppId(String appId) {
        this.appId = appId;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public void setSpeaker(String speaker) {
        this.speaker = speaker;
    }
}