# 常见问题

## 1. 如何配置大模型的各个参数？

以下参数只做为示例，请根据实际情况进行修改。

例如：配置星火大模型 **其他配置** 参数例如：


```
appId=cbd1787
apiSecret=MWUzMzVjNmFhYjNkNDQ5N2E4OWVlpzQz
```

![how_config_model.png](resource/how_config_model.png)


## 2. Bots中调用大模型为空？

![not_find_model_ques.png](resource/not_find_model_ques.png)

因为您没有配置大模型，您可以在[模型管理](http://localhost:8899/ai/llms)中查看大模型列表，并点击菜单 **大模型** 进行配置。在新增大模型的时候

![not_find_model_ques_first.png](resource/not_find_model_ques_first.png)

记得勾选上 **对话模型** 才能在Bots中选择并且调用大模型。

![not_find_model_ques_second.png](resource/not_find_model_ques_second.png)


