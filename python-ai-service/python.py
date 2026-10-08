# 导入 FastAPI 框架，用来搭建 HTTP 接口服务
from fastapi import FastAPI

# 导入类型注解，接收前端传来的 JSON 字典
from typing import Dict

# 导入我们自己写的 RAG 向量检索类
from rag import FAISSRAG

# 导入大模型回答生成函数（智谱AI）
from LLM import generate_answer

# 初始化 FastAPI 应用
app = FastAPI()

# 🔥 初始化 RAG 实例
# 现在不会自动加载本地文件了！
# 完全由 Go 上传同步过来
rag = FAISSRAG()

# ------------------------------
# 测试接口：测试服务是否正常运行
# ------------------------------
@app.post("/test")
def test(data: Dict):
    return {"status": "ok", "msg": "服务正常"}

# ------------------------------
# 同步接口（Go 自动调用）
# ------------------------------
@app.post("/sync")
def sync(data: Dict):
    # 如果没有传 content 内容，就返回错误
    if not data.get("content"):
        return {"status": "error", "msg": "content 不能为空"}
    
    # 把 Go 传来的文档加入向量库
    rag.add_documents([data])
    
    # 返回成功
    return {"status": "success"}

# ===============================
# 🔥 删除向量接口（Go 自动调用）
# ===============================
@app.post("/delete")
def delete_doc(data: Dict):
    doc_id = data.get("id")
    if not doc_id:
        return {"status": "error", "msg": "id 不能为空"}

    # 删除对应 id 的文档和向量
    rag.delete_by_id(str(doc_id))
    return {"status": "success", "msg": "删除成功"}

# ------------------------------
# 🔥 核心聊天接口
# 用户提问 → 检索知识库 → 调用大模型 → 返回答案
# ------------------------------
@app.post("/chat")
def chat(data: Dict):
    # 1. 从前端传来的 JSON 中取出问题
    question = data.get("question", "")
    
    # 2. 如果问题为空，返回提示
    if not question:
        return {"answer": "请输入问题"}
    
    # 3. 🔥 调用 RAG 检索：从知识库找最相关的内容
    hits = rag.search(question, top_k=1)
    
    # 4. 如果没找到匹配内容，返回提示
    if not hits:
        return {"answer": "抱歉，暂无匹配答案"}
    
    # 5. 🔥 调用大模型（智谱）生成最终回答
    answer = generate_answer(question, hits)
    
    # 6. 返回答案给前端/Go服务
    return {"answer": answer, "hits": hits}