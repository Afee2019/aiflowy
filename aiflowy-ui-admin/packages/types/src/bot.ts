interface BotInfo {
  alias: string;
  anonymousEnabled: boolean;
  created: string;
  createdBy: number;
  deptId: number;
  description: string;
  icon: string;
  id: string;
  llmId: string;
  llmOptions: {
    maxMessageCount: number;
    maxReplyLength: number;
    systemPrompt?: string;
    temperature: number;
    topK: number;
    topP: number;
    welcomeMessage: string;
  };
  modified: string;
  modifiedBy: number;
  options: {
    anonymousEnabled: boolean;
    enableDeepThinking: boolean;
    EncodingAESKey: string;
    presetQuestions: presetQuestionsType[];
    reActModeEnabled: boolean;
    voiceEnabled: boolean;
    weChatMpAesKey: string;
    weChatMpAppId: string;
    weChatMpSecret: string;
    weChatMpToken: string;
    welcomeMessage?: string;
  };
  tenantId: number;
  title: string;
  categoryId: any;
  status: number;
}

interface Session {
  botId: string;
  sessionId: string;
  title: string;
}

interface presetQuestionsType {
  key: string;
  description: string;
}

interface ChatMessage {
  role: 'assistant' | 'user';
  created: number;
  updateAt: number;
  id: string;
  content: string;
  options?: {
    type: number;
    user_input: string;
  };
  placement: 'end' | 'start';
  loading?: boolean;
  typing?: boolean;
  fileList?: File[];
  isFog?: boolean;
}

interface AiLlm {
  brand: string;
  deptId: number;
  description: string;
  icon: string;
  id: string;
  llmApiKey: string;
  llmEndpoint: string;
  llmModel: string;
  options: {
    chatPath: string;
    embedPath: string;
    llmEndpoint: string;
    rerankPath: string;
  };
  supportAudioToAudio: boolean;
  supportChat: boolean;
  supportEmbed: boolean;
  supportFeatures: string[];
  supportFunctionCalling: boolean;
  supportImageToImage: boolean;
  supportImageToVideo: boolean;
  supportReranker: boolean;
  supportTextToAudio: boolean;
  supportTextToImage: boolean;
  supportTextToVideo: boolean;
  tenantId: number;
  title: string;
}

export type { AiLlm, BotInfo, ChatMessage, Session };
