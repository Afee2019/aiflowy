# iframe 嵌入

## 不携带 token 嵌入

不携带 token， 客户端将不展现任何后台聊天记录

参数解析: `http://localhost:8899/ai/externalBot/281900928474001408?isIframe=true` <br/>
访问的 Bot 的id : 281900928474001408<br/>
前端网页地址: `http://localhost:8899/ai/externalBot/` <br/>
需要携带固定参数 **isIframe=true**
请求示例如下：
```js
<iframe id="testFrame" src="http://localhost:8899/ai/externalBot/281900928474001408?isIframe=true"></iframe>
```

## 携带 token 嵌入
携带 token， 客户端将展现该用户的后台聊天记录
authKey: 需要后台管理页面生成访问页面的 token
请求示例如下：
```js
<iframe id="testFrame" src="http://localhost:8899/ai/externalBot/281900928474001408?authKey=xxxxxx"></iframe>
```
