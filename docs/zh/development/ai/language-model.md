# 大语言模型

大语言模型在 **AIFlowy** 中（严格意义上是在 **AgentsFlex** 中）是通过 **Llm.java** 来定义的。**Llm.java** 用于定义和描述了大语言模型的基本能力。

**Llm** 针对不同的大语言模型又会有不同的实现类，比如：

- <span style=" color: black;">OpenAiLlm</span> 用于对 ChatGPT 的封装
- <span style=" color: black;">LlamaLlm</span> 用于对 Llama 的封装
- <span style=" color: black;">ChatglmLlm</span> 用于对智谱 AI 的封装
- <span style=" color: black;">SparkLlm</span> 用于对星火大模型的封装
- <span style=" color: black;">QwenLlm</span> 用于对阿里云通义千问的封装

## 示例代码

### 基本对话

```java
OpenAiLlmConfig config = new OpenAiLlmConfig();
config.setApiKey("sk-rts5NF6n*******");

Llm llm = new OpenAiLlm(config);
String response = llm.chat("请问你叫什么名字");
System.out.println(response);
```

### Embedding

```java
Llm llm = OpenAiLlm.of("sk-rts5NF6n*******");
VectorData embeddings = llm.embed(Document.of("some document text"));
System.out.println(Arrays.toString(embeddings.getVector()));
```

### Function Calling

```java
OpenAiLlmConfig config = new OpenAiLlmConfig();
config.setApiKey("sk-rts5NF6n*******");

OpenAiLlm llm = new OpenAiLlm(config);

FunctionPrompt prompt = new FunctionPrompt("今天北京的天气怎么样", WeatherFunctions.class);
FunctionMessageResponse response = llm.chat(prompt);

Object result = response.invoke();

System.out.println(result);
// "Today it will be dull and overcast in 北京"
```