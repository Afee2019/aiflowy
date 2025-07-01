import React, {useLayoutEffect, useMemo, useRef, useState} from 'react';
import {Bubble, Prompts, Sender, ThoughtChain, ThoughtChainItem, Welcome} from '@ant-design/x';
import {Button, GetProp, message, Space, Spin, Typography} from 'antd';
import {CopyOutlined, FolderAddOutlined, OpenAIOutlined, SyncOutlined} from '@ant-design/icons';
// import ReactMarkdown from 'react-markdown';
// import remarkGfm from 'remark-gfm';
// import remarkBreaks from 'remark-breaks';
import logo from "/favicon.png";
import {UserOutlined} from '@ant-design/icons';
import './aiprochat.less'
import markdownit from 'markdown-it';

const fooAvatar: React.CSSProperties = {
    color: '#fff',
    backgroundColor: '#87d068',
};

export type ChatMessage = {
    id: string;
    content: string;
    role: 'user' | 'assistant' | 'aiLoading' | string;
    created: number;
    updateAt?: number;
    loading?: boolean;
    thoughtChains?: Array<ThoughtChainItem>
    options?: object;
};

// 事件类型
export type EventType = 'thinking' | 'thought' | 'toolCalling' | 'callResult' | string;

export type EventHandlerResult = {
    handled: boolean; // 是否已处理该事件
    data?: any; // 处理结果数据
};

// 事件处理器函数类型
export type EventHandler = (eventType: EventType, eventData: any, context: {
    chats: ChatMessage[];
    setChats: (value: ((prevState: ChatMessage[]) => ChatMessage[]) | ChatMessage[]) => void;
}) => EventHandlerResult | Promise<EventHandlerResult>;


export type AiProChatProps = {
    loading?: boolean;
    chats?: ChatMessage[];
    onChatsChange?: (value: ((prevState: ChatMessage[]) => ChatMessage[]) | ChatMessage[]) => void;
    style?: React.CSSProperties;
    appStyle?: React.CSSProperties;
    helloMessage?: string;
    botAvatar?: string;
    request: (messages: ChatMessage[]) => Promise<Response>;
    clearMessage?: () => void;
    showQaButton?: boolean;
    onQaButtonClick?: (currentChat: ChatMessage, index: number, allChats: ChatMessage[]) => void;
    prompts?: GetProp<typeof Prompts, 'items'>;
    inputDisabled?: boolean;
    customToolBarr?: React.ReactNode;
    onCustomEvent?: EventHandler;
    onCustomEventComplete?: EventHandler;
};

export const RenderMarkdown: React.FC<{ content: string }> = ({content}) => {

    const md = markdownit({html: true, breaks: true});

    return (
        <Typography>
            {/* biome-ignore lint/security/noDangerouslySetInnerHtml: used in demo */}
            <div dangerouslySetInnerHTML={{__html: md.render(content)}}/>
        </Typography>
    );
};

export const AiProChat = ({
                              loading,
                              chats: parentChats,
                              onChatsChange: parentOnChatsChange,
                              style = {},
                              appStyle = {},
                              helloMessage = '欢迎使用 AIFlowy',
                              botAvatar = `${logo}`,
                              request,
                              showQaButton = false,
                              onQaButtonClick = (): void => {
                              },
                              clearMessage,
                              inputDisabled = false,
                              prompts,
                              customToolBarr,
                              onCustomEvent,
                              onCustomEventComplete,
                          }: AiProChatProps) => {
    const isControlled = parentChats !== undefined && parentOnChatsChange !== undefined;
    const [internalChats, setInternalChats] = useState<ChatMessage[]>([]);
    const chats = useMemo(() => {
        return isControlled ? parentChats : internalChats;
    }, [isControlled, parentChats, internalChats]);
    const setChats = isControlled ? parentOnChatsChange : setInternalChats;
    const [content, setContent] = useState('');
    const [sendLoading, setSendLoading] = useState(false);
    const [isStreaming, setIsStreaming] = useState(false);
    const messagesContainerRef = useRef<HTMLDivElement>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    // 控制是否允许自动滚动
    const autoScrollEnabled = useRef(true); // 默认允许自动滚动
    const isUserScrolledUp = useRef(false); // 用户是否向上滚动过

    //  使用 ref 来跟踪事件状态，避免异步状态更新问题
    const currentEventType = useRef<string | null>(null);
    const eventContent = useRef<string>(''); // 当前事件累积的内容

    useRef<string | null>(null);
    // 滚动到底部逻辑
    const scrollToBottom = () => {
        const container = messagesContainerRef.current;
        if (container && autoScrollEnabled.current) {
            container.scrollTop = container.scrollHeight;
        }
    };

    // 组件挂载时滚动
    useLayoutEffect(() => {
        scrollToBottom();
    }, []);

    // 消息更新时滚动
    useLayoutEffect(() => {
        if (autoScrollEnabled.current) {
            scrollToBottom();
        }
    }, [chats]);
    useLayoutEffect(() => {
        const container = messagesContainerRef.current;
        if (!container) return;

        const handleScroll = () => {
            const {scrollTop, scrollHeight, clientHeight} = container;
            const atBottom = scrollHeight - scrollTop <= clientHeight + 5; // 允许误差 5px

            if (atBottom) {
                // 用户回到底部，恢复自动滚动
                autoScrollEnabled.current = true;
                isUserScrolledUp.current = false;
            } else {
                // 用户向上滚动，禁用自动滚动
                autoScrollEnabled.current = false;
                isUserScrolledUp.current = true;
            }
        };

        container.addEventListener('scroll', handleScroll);
        return () => {
            container.removeEventListener('scroll', handleScroll);
        };
    }, []);

    // 处理事件进度（事件进行中）
    const handleEventProgress = async (eventType: EventType, eventData: any): Promise<boolean> => {
        if (onCustomEvent) {
            try {

                const result = await onCustomEvent(eventType, eventData, {
                    chats,
                    setChats,
                });

                if (result.handled) {
                    console.log(`Event progress "${eventType}" handled by custom handler`);
                    return true;
                }
            } catch (error) {
                console.error(`Custom event progress handler error for "${eventType}":`, error);
            }
        }


        // 使用现有的默认处理逻辑
        return handleDefaultEvent(eventType, eventData);
    };

    // 处理事件完成
    const handleEventComplete = async (eventType: EventType, finalContent: string): Promise<boolean> => {

        const eventData = {
            content: finalContent,
            accumulatedContent: finalContent,
            isComplete: true
        };

        if (onCustomEventComplete) {
            try {
                const result = await onCustomEventComplete(eventType, eventData, {
                    chats,
                    setChats
                });

                if (result.handled) {
                    console.log(`Event complete "${eventType}" handled by custom complete handler`);
                    return true;
                }
            } catch (error) {
                console.error(`Custom event complete handler error for "${eventType}":`, error);
            }
        }


        // 使用现有的默认处理逻辑
        return handleDefaultEvent(eventType, eventData);
    };


    const handleDefaultEvent = (eventType: EventType, eventData: any): boolean => {

        if (eventData.isComplete || eventType === "content") {
            return true;
        }

        // 🧠 处理 ThoughtChain 相关事件
        if (['thinking', 'thought', 'toolCalling', 'callResult'].includes(eventType)) {

            setChats((prevChats: ChatMessage[]) => {
                const newChats = [...prevChats];

                const lastAiIndex = (() => {
                    for (let i = newChats.length - 1; i >= 0; i--) {
                        if (newChats[i].role === 'assistant') {
                            return i;
                        }
                    }
                    return -1;
                })();

                const aiMessage = newChats[lastAiIndex];
                aiMessage.loading = false;

                return newChats;
            });

            setChats((prevChats: ChatMessage[]) => {
                const newChats = [...prevChats];

                // 找到最后一条 assistant 消息
                const lastAiIndex = (() => {
                    for (let i = newChats.length - 1; i >= 0; i--) {
                        if (newChats[i].role === 'assistant') {
                            return i;
                        }
                    }
                    return -1;
                })();

                if (lastAiIndex !== -1) {
                    const aiMessage = newChats[lastAiIndex];

                    // 初始化 thoughtChains 数组（如果不存在）
                    if (!aiMessage.thoughtChains) {
                        aiMessage.thoughtChains = [];
                    }

                    const title = eventData.metadataMap.chainTitle;
                    const description = (eventData.accumulatedContent || eventData.content || '') as string;

                    // 获取事件ID
                    const eventId = eventData.id || eventData.metadataMap?.id;

                    if (eventId) {
                        // 查找是否存在相同 id 的思维链项
                        const targetIndex = aiMessage.thoughtChains.findIndex(item =>
                            item.key === eventId || item.key === String(eventId)
                        );

                        if (targetIndex !== -1) {
                            // 找到相同 id 的项，更新该项
                            aiMessage.thoughtChains[targetIndex] = {
                                ...aiMessage.thoughtChains[targetIndex],
                                key: eventId,
                                title,
                                content: <RenderMarkdown content={description} />,
                                status: 'pending'
                            };
                            console.log(`Updated ThoughtChain item with id: ${eventId} for event: ${eventType}`);
                        } else {
                            // 没找到相同 id 的项，创建新项
                            const newItem: ThoughtChainItem = {
                                key: eventId,
                                title,
                                content:  <RenderMarkdown content={description} />,
                                status: 'pending'
                            };

                            aiMessage.thoughtChains.push(newItem);
                            console.log(`Created new ThoughtChain item with id: ${eventId} for event: ${eventType}`);
                        }
                    } else {
                        console.warn(`Event ${eventType} has no id, skipping ThoughtChain processing`);
                    }

                    // 更新消息的更新时间
                    aiMessage.updateAt = Date.now();
                }

                return newChats;
            });

            return true;
        }

        return true;
    };

    // 提交流程优化
    const handleSubmit = async (newMessage: string) => {
        const messageContent = newMessage?.trim() || content.trim();
        if (!messageContent) return;

        setSendLoading(true);
        setIsStreaming(true);

        const userMessage: ChatMessage = {
            role: 'user',
            id: Date.now().toString(),
            content: messageContent,
            created: Date.now(),
            updateAt: Date.now(),
        };

        const aiMessage: ChatMessage = {
            role: 'assistant',
            id: Date.now().toString(),
            content: '',
            loading: true,
            created: Date.now(),
            updateAt: Date.now(),
        };

        const temp = [userMessage, aiMessage];
        setChats?.((prev: ChatMessage[]) => [...(prev || []), ...temp]);
        setTimeout(scrollToBottom, 50);
        setContent('');

        try {
            const response = await request([...(chats || []), userMessage]);
            if (!response?.body) return;

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let partial = '';
            let currentContent = '';
            let typingIntervalId: NodeJS.Timeout | null = null;

            // 用于等待打字效果完成的Promise
            const waitForTypingComplete = (): Promise<void> => {
                return new Promise((resolve) => {
                    const checkTypingComplete = () => {
                        if (currentContent === partial) {
                            resolve();
                        } else {
                            setTimeout(checkTypingComplete, 50);
                        }
                    };
                    checkTypingComplete();
                });
            };

            let isStreamFinished = false;
            let shouldContinueReading = true;
            //  重置事件状态
            currentEventType.current = null;
            eventContent.current = '';

            while (shouldContinueReading) {
                const {done, value} = await reader.read();
                if (done) {
                    isStreamFinished = true;
                    shouldContinueReading = false;
                    //  流结束时，如果还有未完成的事件，触发事件完成处理
                    if (currentEventType.current) {
                        console.log(`Stream finished, completing event: ${currentEventType.current}`);
                        await handleEventComplete(currentEventType.current, eventContent.current);
                        currentEventType.current = null;
                        eventContent.current = '';
                    }
                    break;
                }

                const decode = decoder.decode(value, {stream: true});
                const parse = JSON.parse(decode);
                const respData = JSON.parse(parse.data);

                const incomingEventType = parse.event || 'content';

                // 检查是否切换到了新的事件类型（使用 ref.current）
                if (currentEventType.current && currentEventType.current !== incomingEventType) {
                    console.log(`Event type changed from ${currentEventType.current} to ${incomingEventType}, completing previous event`);

                    try {
                        // 上一个事件完成，触发完成处理
                        await handleEventComplete(currentEventType.current, eventContent.current);
                    } catch (error) {
                        console.error(` Event transition failed:`, error);
                    }

                    // 重置累积内容
                    eventContent.current = '';
                }

                //  更新当前事件类型
                currentEventType.current = incomingEventType;

                if (incomingEventType !== 'content') {
                    // 累积事件内容
                    const newEventContent = eventContent.current + (respData.content || '');
                    eventContent.current = newEventContent;

                    try {
                        //  事件处理失败时直接抛出错误
                        const eventHandled = await handleEventProgress(incomingEventType, {
                            ...respData,
                            accumulatedContent: newEventContent,
                            isComplete: false
                        });

                        // 如果事件已被处理，跳过内容更新逻辑
                        if (eventHandled) {
                            continue;
                        }
                    } catch (error) {
                        console.error(`Event processing failed, terminating stream:`, error);
                    }
                }

                // 处理内容更新
                partial += respData.content || '';

                // 清除之前的打字间隔
                if (typingIntervalId) {
                    clearInterval(typingIntervalId);
                }

                // 开始新的打字效果
                typingIntervalId = setInterval(() => {
                    if (currentContent.length < partial.length) {
                        currentContent = isStreamFinished ? partial : partial.slice(0, currentContent.length + 2);
                        setChats?.((prev: ChatMessage[]) => {
                            const newChats = [...(prev || [])];
                            const lastMsg = newChats[newChats.length - 1];
                            if (!lastMsg) return prev;

                            if (lastMsg?.role === 'assistant') {
                                lastMsg.loading = false;
                                lastMsg.content = currentContent;
                                lastMsg.updateAt = Date.now();
                            }
                            return newChats;
                        });

                        if (autoScrollEnabled.current) {
                            scrollToBottom();
                        }
                    }

                    // 当前内容已经追上完整内容时停止
                    if (currentContent == partial || isStreamFinished) {
                        clearInterval(typingIntervalId!);
                        typingIntervalId = null;
                    }
                }, 50);
            }

            // 等待最后的打字效果完成
            await waitForTypingComplete();

            // 清理间隔（如果还存在）
            if (typingIntervalId) {
                clearInterval(typingIntervalId);
            }

            setChats((prev: ChatMessage[]) => {
                console.log(prev);
                const newChats = prev;
                if (prev){
                    const chatMessage = newChats[prev.length - 1];
                    if (chatMessage){
                        chatMessage.content?.replace("Final Answer:","");
                    }
                }
                return newChats;
            })

        } catch (error) {
            console.error(`Stream processing error:`, error);
        } finally {
            // 确保打字效果完成后再重置状态
            setIsStreaming(false);
            setSendLoading(false);
            console.log(chats)
        }
    };

    // 重新生成消息
    const handleRegenerate = async (index: number) => {
        // 找到当前 assistant 消息对应的上一条用户消息
        const prevMessage: ChatMessage = {
            role: 'user',
            id: Date.now().toString(),
            content: chats[index - 1].content,
            loading: false,
            created: Date.now(),
            updateAt: Date.now(),
        };
        setContent(prevMessage.content)
        const aiMessage: ChatMessage = {
            role: 'assistant',
            id: Date.now().toString(),
            content: '',
            loading: true,
            created: Date.now(),
            updateAt: Date.now(),
        };
        setSendLoading(true);
        setIsStreaming(true);
        const temp = [prevMessage, aiMessage];
        setChats?.((prev: ChatMessage[]) => [...(prev || []), ...temp]);
        setTimeout(scrollToBottom, 50);
        setContent('');

        try {
            const response = await request([...(chats || []), prevMessage]);
            if (!response?.body) return;

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let partial = '';
            let currentContent = '';
            let typingIntervalId: NodeJS.Timeout | null = null;

            // 用于等待打字效果完成的Promise
            const waitForTypingComplete = (): Promise<void> => {
                return new Promise((resolve) => {
                    const checkTypingComplete = () => {
                        if (currentContent === partial) {
                            resolve();
                        } else {
                            setTimeout(checkTypingComplete, 50);
                        }
                    };
                    checkTypingComplete();
                });
            };

            let isStreamFinished = false;
            let shouldContinueReading = true;

            //  重置事件状态
            currentEventType.current = null;
            eventContent.current = '';

            while (shouldContinueReading) {
                const {done, value} = await reader.read();
                if (done) {
                    isStreamFinished = true;
                    shouldContinueReading = false;

                    //  流结束时，如果还有未完成的事件，触发事件完成处理
                    if (currentEventType.current) {
                        console.log(`Regenerate stream finished, completing event: ${currentEventType.current}`);
                        await handleEventComplete(currentEventType.current, eventContent.current);
                        currentEventType.current = null;
                        eventContent.current = '';
                    }
                    continue;
                }

                const decode = decoder.decode(value, {stream: true});

                //  检查是否为包含事件的格式
                try {
                    const parse = JSON.parse(decode);
                    const respData = JSON.parse(parse.data);
                    const incomingEventType = parse.event || 'content';

                    //  检查是否切换到了新的事件类型
                    if (currentEventType.current && currentEventType.current !== incomingEventType) {
                        console.log(`Regenerate event type changed from ${currentEventType.current} to ${incomingEventType}, completing previous event`);

                        // 上一个事件完成，触发完成处理
                        await handleEventComplete(currentEventType.current, eventContent.current);

                        // 重置累积内容
                        eventContent.current = '';
                    }

                    //  更新当前事件类型
                    currentEventType.current = incomingEventType;

                    if (incomingEventType !== 'content') {
                        //  累积事件内容
                        const newEventContent = eventContent.current + (respData.content || '');
                        eventContent.current = newEventContent;

                        //  处理事件进度
                        const eventHandled = await handleEventProgress(incomingEventType, {
                            ...respData,
                            accumulatedContent: newEventContent,
                            isComplete: false
                        });

                        // 如果事件已被处理，跳过内容更新逻辑
                        if (eventHandled) {
                            continue;
                        }
                    }

                    // 处理内容更新
                    partial += respData.content || '';
                } catch (error) {
                    //  如果解析失败，当作普通内容处理（兼容旧格式）
                    partial += decode;
                }

                // 清除之前的打字间隔
                if (typingIntervalId) {
                    clearInterval(typingIntervalId);
                }

                // 开始新的打字效果
                typingIntervalId = setInterval(() => {
                    if (currentContent.length < partial.length) {
                        currentContent = isStreamFinished ? partial : partial.slice(0, currentContent.length + 2);
                        setChats?.((prev: ChatMessage[]) => {
                            const newChats = [...(prev || [])];
                            const lastMsg = newChats[newChats.length - 1];

                            if (!lastMsg) {
                                return prev;
                            }

                            if (lastMsg.role === 'assistant') {
                                lastMsg.loading = false;
                                lastMsg.content = currentContent;
                                lastMsg.updateAt = Date.now();
                            }
                            return newChats;
                        });

                        if (autoScrollEnabled.current) {
                            scrollToBottom();
                        }
                    }

                    // 当前内容已经追上完整内容时停止
                    if (currentContent === partial || isStreamFinished) {
                        clearInterval(typingIntervalId!);
                        typingIntervalId = null;
                    }
                }, 50);
            }

            // 等待最后的打字效果完成
            await waitForTypingComplete();

            // 清理间隔（如果还存在）
            if (typingIntervalId) {
                clearInterval(typingIntervalId);
            }

        } catch (error) {
            console.error('Regenerate error:', error);
        } finally {
            // 确保打字效果完成后再重置状态
            setIsStreaming(false);
            setSendLoading(false);
        }
    };


    // 渲染消息列表
    const renderMessages = () => {
        if (!chats?.length) {
            return (
                <Welcome
                    variant="borderless"
                    icon={<img
                        src={botAvatar}
                        style={{width: 32, height: 32, borderRadius: '50%'}}
                        alt="AI Avatar"
                    />}
                    description={helloMessage}
                    styles={{icon: {width: 40, height: 40}}}
                />
            );
        }
        return (
            <Bubble.List
                autoScroll={true}
                items={chats.map((chat, index) => ({
                    key: chat.id + Math.random().toString(),
                    // typing: {suffix: <>💗</>},
                    header: (
                        <Space>
                            {new Date(chat.created).toLocaleString()}
                        </Space>
                    ),
                    loading: chat.loading,
                    loadingRender: () => (
                        <Space>
                            <Spin size="small"/>
                            AI正在思考中...
                        </Space>
                    ),
                    footer: (
                        <Space>
                            {(chat.role === 'assistant') && !isStreaming && (<Button
                                color="default"
                                variant="text"
                                size="small"
                                icon={<SyncOutlined/>}
                                onClick={() => {
                                    // 点击按钮时重新生成该消息
                                    if (chat.role === 'assistant') {
                                        handleRegenerate(index);
                                    }
                                }}
                            />)}

                            {
                                !isStreaming && <Button
                                    color="default"
                                    variant="text"
                                    size="small"
                                    icon={<CopyOutlined/>}
                                    onClick={async () => {
                                        try {
                                            await navigator.clipboard.writeText(chat.content);
                                            message.success('复制成功');
                                        } catch (error) {
                                            console.log(error);
                                            message.error('复制失败');
                                        }
                                    }}
                                />
                            }
                            {(chat.role === 'user' && showQaButton) && !isStreaming && <Button
                                color="default"
                                variant="text"
                                size="small"

                                icon={<FolderAddOutlined/>}
                                onClick={async () => {
                                    handleQaClick(chat, index)
                                }}
                            ></Button>}
                        </Space>
                    ),
                    role: chat.role === 'user' ? 'local' : 'ai',
                    content: chat.role === 'assistant' ? (
                        <div>
                            {/* 🧠 使用 ThoughtChain 组件 */}
                            {chat.thoughtChains && chat.thoughtChains.length > 0 && (
                                <ThoughtChain
                                    items={chat.thoughtChains}
                                    style={{marginBottom: '12px'}}
                                />
                            )}

                            {/* 🌟 渲染主要内容 */}
                            <RenderMarkdown content={chat.content} />
                        </div>
                    ) : chat.content,
                    avatar: chat.role === 'assistant' ? (
                        <img
                            src={botAvatar}
                            style={{width: 32, height: 32, borderRadius: '50%'}}
                            alt="AI Avatar"
                        />
                    ) : {icon: <UserOutlined/>, style: fooAvatar},
                }))}
                roles={{ai: {placement: 'start'}, local: {placement: 'end'}}}
            />
        );
    };

    // qa按钮点击事件
    const handleQaClick = (chat: ChatMessage, index: number) => {
        if (onQaButtonClick) {
            onQaButtonClick(chat, index, chats);
        }
    };

    const SENDER_PROMPTS = prompts || [
        {
            key: '1',
            description: '你好'
        },
        {
            key: '2',
            description: '你是谁？'
        }
    ];


    return (
        <div
            style={{
                width: '100%',
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                background: '#fff',
                border: '1px solid #f3f3f3',
                ...appStyle,
                ...style,
            }}
        >
            {/* 消息容器 */}
            <div
                ref={messagesContainerRef}
                style={{
                    flex: 1,
                    overflowY: 'auto',
                    padding: '16px',
                    scrollbarWidth: 'none',
                }}
            >
                {loading ? (
                    <Spin tip="加载中..."/>
                ) : (
                    <>
                        {renderMessages()}
                        <div ref={messagesEndRef}/>
                        {/* 锚点元素 */}
                    </>
                )}
            </div>
            {/* 输入区域 */}

            <div
                style={{
                    borderTop: '1px solid #eee',
                    padding: '12px',
                    display: 'flex',
                    flexDirection: "column",
                    gap: '8px',
                }}
            >

                {/* 🌟 提示词 */}
                <Prompts
                    items={SENDER_PROMPTS}
                    onItemClick={(info) => {
                        handleSubmit(info.data.description as string)
                    }}
                    styles={{
                        item: {padding: '6px 12px'},
                    }}
                />

                {customToolBarr ?
                    <div style={{
                        width: "100%",
                        display: "flex",
                        justifyContent: "start",
                        alignItems: "center",
                    }}>
                        {customToolBarr}
                    </div> : <></>
                }

                <Sender
                    value={content}
                    onChange={setContent}
                    onSubmit={handleSubmit}
                    loading={sendLoading || isStreaming}
                    disabled={inputDisabled}
                    actions={(_, info) => (
                        <Space size="small">
                            <info.components.ClearButton
                                disabled={sendLoading || isStreaming || !chats?.length}  // 强制不禁用
                                title="删除对话记录"
                                style={{fontSize: 20}}
                                onClick={async (e) => {
                                    e.preventDefault();  // 阻止默认行为（如果有）
                                    setSendLoading(true)
                                    await clearMessage?.();
                                    setSendLoading(false)
                                }}
                            />
                            <info.components.SendButton
                                type="primary"
                                disabled={inputDisabled}
                                icon={<OpenAIOutlined/>}
                                loading={sendLoading || isStreaming}
                            />
                        </Space>
                    )}
                />
            </div>
        </div>
    );
};