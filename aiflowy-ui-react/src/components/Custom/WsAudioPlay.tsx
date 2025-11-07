import React, {forwardRef, useCallback, useEffect, useImperativeHandle, useRef, useState} from "react"

export type WsAudioPlayProps = {
    ref?: any;
    sessionId: any;
    playStateChange?: (isPlaying: boolean) => void;
}
const authKey = `${import.meta.env.VITE_APP_AUTH_KEY || "authKey"}`;
export const WsAudioPlay: React.FC<WsAudioPlayProps> = forwardRef((props, ref) => {

    const wsRef = useRef<WebSocket | null>(null)
    const voiceMapRef = useRef<Map<string, string[]>>(new Map());
    const audioRef = useRef<HTMLAudioElement>(null);
    const mediaSourceRef = useRef<MediaSource | null>(null);
    const sourceBufferRef = useRef<SourceBuffer | null>(null);
    const currentMessageIdRef = useRef<string | null>(null);
    const isInitializedRef = useRef(false);
    const audioQueueRef = useRef<ArrayBuffer[]>([]); // 音频数据队列
    const isProcessingRef = useRef(false); // 是否正在处理队列
    const [isPlaying, setIsPlaying] = useState(false);

    const {sessionId, playStateChange} = props

    useEffect(() => {
        playStateChange?.(isPlaying);
    }, [isPlaying, playStateChange]);

    const token = localStorage.getItem(authKey)

    // 发送文本到WebSocket
    const sendText = useCallback((text: string) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            wsRef.current.send(text);
            //console.log('Text sent to server:', text);
        }
    }, [sessionId]);

    // 初始化MediaSource
    const initMediaSource = useCallback(() => {
        if (mediaSourceRef.current) {
            if (mediaSourceRef.current.readyState === 'open') {
                mediaSourceRef.current.endOfStream();
            }
            mediaSourceRef.current = null;
        }

        // 清空队列
        audioQueueRef.current = [];
        isProcessingRef.current = false;

        const mediaSource = new MediaSource();
        mediaSourceRef.current = mediaSource;
        isInitializedRef.current = false;

        mediaSource.addEventListener('sourceopen', () => {
            //console.log('MediaSource opened');
            try {
                const sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg');
                sourceBufferRef.current = sourceBuffer;
                isInitializedRef.current = true;

                sourceBuffer.addEventListener('updateend', () => {
                    //console.log('SourceBuffer update ended');
                    isProcessingRef.current = false;
                    // 检查队列中是否有更多数据需要处理
                    processQueue();
                });

                sourceBuffer.addEventListener('error', (e) => {
                    console.error('SourceBuffer error:', e);
                    isProcessingRef.current = false;
                });

            } catch (err) {
                console.error('Error creating SourceBuffer:', err);
            }
        });

        mediaSource.addEventListener('error', (e) => {
            console.error('MediaSource error:', e);
        });

        if (audioRef.current) {
            audioRef.current.src = URL.createObjectURL(mediaSource);
        }
    }, []);

    // 处理音频队列
    const processQueue = useCallback(() => {
        if (!sourceBufferRef.current || !isInitializedRef.current) {
            return;
        }

        // 如果正在处理或者队列为空，直接返回
        if (isProcessingRef.current || audioQueueRef.current.length === 0) {
            return;
        }

        // 如果SourceBuffer正在更新，等待下次处理
        if (sourceBufferRef.current.updating) {
            return;
        }

        try {
            const audioBuffer = audioQueueRef.current.shift();
            if (audioBuffer) {
                isProcessingRef.current = true;
                sourceBufferRef.current.appendBuffer(audioBuffer);
            }
        } catch (err) {
            console.error('Error processing audio queue:', err);
            isProcessingRef.current = false;
        }
    }, []);
    // 停止播放
    const stop = useCallback(() => {
        if (audioRef.current) {
            audioRef.current.pause();
            audioRef.current.currentTime = 0;
            setIsPlaying(false);
        }

        // 清理SourceBuffer
        if (sourceBufferRef.current && !sourceBufferRef.current.updating) {
            try {
                if (sourceBufferRef.current &&
                    mediaSourceRef.current &&
                    mediaSourceRef.current.readyState === 'open' &&
                    !sourceBufferRef.current.updating) {
                    try {
                        sourceBufferRef.current.abort();
                    } catch (err) {
                        console.warn('Error aborting SourceBuffer (non-critical):', err);
                    }
                }
            } catch (err) {
                console.error('Error aborting SourceBuffer:', err);
            }
        }

        // 清空队列
        audioQueueRef.current = [];
        isProcessingRef.current = false;

        currentMessageIdRef.current = null;
    }, []);
    // 添加音频数据到队列
    const appendAudioData = useCallback((base64Data: string) => {
        if (!sourceBufferRef.current || !isInitializedRef.current) {
            console.warn('SourceBuffer not ready');
            return false;
        }

        try {
            // Base64 转 ArrayBuffer
            const binaryString = atob(base64Data);
            const bytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            const audioBuffer = bytes.buffer;

            // 将数据添加到队列
            audioQueueRef.current.push(audioBuffer);

            // 开始处理队列
            processQueue();

            return true;
        } catch (err) {
            console.error('Error appending audio data:', err);
            return false;
        }
    }, [processQueue]);

    // 开始流式播放
    const startStreamPlayback = useCallback((messageId: string) => {
        // 停止当前播放
        stop();

        // 重置MediaSource
        initMediaSource();

        // 设置当前播放的messageId
        currentMessageIdRef.current = messageId;

        //console.log(`Starting stream playback for message: ${messageId}`);

        // 开始播放
        if (audioRef.current && mediaSourceRef.current) {
            // 等待MediaSource就绪后开始播放
            const checkReady = () => {
                if (mediaSourceRef.current?.readyState === 'open' && isInitializedRef.current) {
                    audioRef.current?.play().catch(err => {
                        console.error('Play failed:', err);
                    });
                } else {
                    setTimeout(checkReady, 10);
                }
            };
            checkReady();
        }
    }, [initMediaSource, stop]);


    // 播放音频
    const play = useCallback((messageId: string) => {
        const voiceMap = voiceMapRef.current;
        if (!voiceMap.has(messageId)) {
            console.warn(`No audio data found for messageId: ${messageId}`);
            return;
        }
        // 停止当前播放
        stop();
        try {
            // 重置MediaSource
            initMediaSource();

            // 设置当前播放的messageId
            currentMessageIdRef.current = messageId;

            // 获取该messageId的所有音频数据
            const base64DataArray = voiceMap.get(messageId)!;
            const fullBase64 = base64DataArray.join('');

            // Base64 转 ArrayBuffer
            const binaryString = atob(fullBase64);
            const bytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            const audioBuffer = bytes.buffer;

            // 等待MediaSource就绪后添加数据
            const checkSourceBuffer = () => {
                if (sourceBufferRef.current && mediaSourceRef.current &&
                    mediaSourceRef.current.readyState === 'open') {

                    // 添加数据到SourceBuffer
                    if (!sourceBufferRef.current.updating) {
                        sourceBufferRef.current.appendBuffer(audioBuffer);

                        // 数据添加完成后开始播放
                        sourceBufferRef.current.addEventListener('updateend', () => {
                            if (mediaSourceRef.current) {
                                mediaSourceRef.current.endOfStream();
                            }
                            if (audioRef.current) {
                                audioRef.current.play().catch(err => {
                                    console.error('播放失败:', err);
                                });
                            }
                        }, {once: true});
                    } else {
                        setTimeout(checkSourceBuffer, 10);
                    }
                } else {
                    setTimeout(checkSourceBuffer, 10);
                }
            };

            checkSourceBuffer();

        } catch (err) {
            console.error(`Error playing audio for messageId ${messageId}:`, err);
        }
    }, [initMediaSource]);


    useImperativeHandle(ref, () => ({
        sendText,
        play,
        stop,
        isPlaying
    }), [sendText, play, stop, isPlaying])

    const initWs = useCallback(() => {
        const websocket = new WebSocket(`${import.meta.env.VITE_APP_WS_SERVER_ENDPOINT}/api/v1/aiBot/ws/audio?sessionId=${sessionId}&token=${token}`);
        wsRef.current = websocket;
        websocket.onopen = () => {
            console.log("ws 连接成功");
        };
        websocket.onerror = (event: Event) => {
            console.error("ws 连接错误:", event);
        };
        websocket.onmessage = (event: MessageEvent) => {
            const audioData: { type: string, content: string, messageId: string } = JSON.parse(event.data)
            //console.log("ws 收到数据:", audioData)
            const type = audioData.type
            const messageId = audioData.messageId
            const content = audioData.content

            const voiceMap = voiceMapRef.current;
            if (type === '_error_') {
                console.error('服务端返回错误', audioData)
            }
            if (type === '_start_') {
                // 开始新的音频流
                startStreamPlayback(messageId);
                // 初始化存储
                voiceMap.set(messageId, []);
            }
            if (type === '_data_') {
                if (!voiceMap.has(messageId)) {
                    voiceMap.set(messageId, []);
                }
                voiceMap.get(messageId)?.push(content);
                // 如果是当前正在播放的消息，立即添加到播放器
                if (currentMessageIdRef.current === messageId) {
                    appendAudioData(content);
                }
            }
            if (type === '_end_') {
                //console.log(`Audio stream ended for message: ${messageId}`);
                // 标记流结束
                if (currentMessageIdRef.current === messageId &&
                    mediaSourceRef.current &&
                    mediaSourceRef.current.readyState === 'open') {
                    // 等待队列处理完成后再结束流
                    const checkQueueEmpty = () => {
                        if (audioQueueRef.current.length === 0 && !isProcessingRef.current) {
                            mediaSourceRef.current?.endOfStream();
                        } else {
                            setTimeout(checkQueueEmpty, 10);
                        }
                    };
                    checkQueueEmpty();
                }
            }
        };
        websocket.onclose = (event) => {
            console.log("ws 连接关闭:", event);
        };
    }, [sessionId, token, startStreamPlayback, appendAudioData])

    useEffect(() => {
        initWs()

        return () => {
            if (wsRef.current) {
                wsRef.current.close(1000, '关闭ws连接');
            }
        };
    }, [initWs])

    return (
        <>
            <audio
                ref={audioRef}
                style={{display: 'none'}}
                onEnded={() => setIsPlaying(false)}
                onPlay={() => setIsPlaying(true)}
                onPause={() => setIsPlaying(false)}
            />
        </>
    )
})