# -
基于 Go(Gin)+ Python(FastAPI) 构建跨语言微服务,Go 负责认证与知识库管理,Python 通过 HTTP 提供 RAG 能力; • 实现 RAG 检索链路:文档经智谱 embedding 向量化写入 FAISS,语义检索召回后由 LangChain 组装 Prompt 调用大模型生成答案; • 实现注册登录与角色权限控制,采用 bcrypt 密码加密 + JWT 无状态认证; • 集成 Viper 配置热加载、Zap+lumberjack 日志切割、Snowflake 分布式 ID、令牌桶限流、优雅关闭。
