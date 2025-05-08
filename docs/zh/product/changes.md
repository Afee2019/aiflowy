# AIFlowy 更新记录

## v1.0.1 2025-05-08
- 新增：导入导出工作流. close #IC59WP
- 新增：用户拿到机器人(外部访问地址), 登录后返回机器人聊天页面 close #IC4MU6
- 新增：机器人增加自定义欢迎语 close #IC560Y
- 优化：优化博查搜索实现
- 优化：加入大模型按钮增加llmModel、brand参数。改善体验
- 优化：升级 Tinyflow 到最新版本
- 优化：优化 bots 的进入方式 close #IC568Q
- 优化：统一弹窗风格. close #IC4NQZ
- 优化：优化外部聊天地址页签和图标，欢迎语
- 优化：增加向量数据库配置说明
- 优化：优化 TinyFlowConfigService 的初始化方法及其位置
- 优化：外部访问 bot 增加页签图标显示和 bot 名称
- 优化：优化 bots 卡片高度，描述内容超出使用 tip 友好提示 close #IC5PKP
- 优化：统一 RAG 检索的相识度获取
- 优化：机器人聊天解除工作流注释 close #IC5J3C
- 修复：通过升级 Agents-Flex 和 Tinyflow 已修复 issues close #IC64VB close #IC60BS close #IC5W27 close #IC5451
- 修复：修复工作流和大模型配置调用显示超时问题 close #IC30BW
- 修复：修复知识库使用 milvus 向量模型为空问题 和 milvus 配置问题 close #IC00HI
- 修复：修复 okhttp 相关依赖版本冲突的问题
- 修复：工作流试运行时先保存 
- 修复：搜索引擎节点添加博查搜索
- 修复：提取设置tinyflow的公共方法
- 修复：修复知识库 redis 作为向量库时返回 Score 结果错误
- 修复：本地模型默认请求地址
- 文档：更新 openSearch 配置说明
- 文档：更新 elasticsearch 配置说明
- 文档：更新插件开发文档
- 文档：工作流自定义节点文档
- 文档：文件内容提取，文件生成节点的文档。
- 文档：新增搜索引擎节点文档.
- 文档：新增内容模板节点文档


## v1.0.0 2025-04-30
完善 Bot 应用、Bot 插件、RAG 知识库、AI 工作流、用户管理、部门管理、角色管理、菜单管理等基础功能