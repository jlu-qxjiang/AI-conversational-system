import os
from typing import List, Dict
import requests
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_core.runnables import RunnableLambda

ZHIPU_API_KEY = os.environ.get("ZHIPU_API_KEY", "")
ZHIPU_API_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"


def _call_zhipu(messages: list) -> str:
    """LangChain Runnable 包装：调用智谱 API"""
    headers = {
        "Authorization": f"Bearer {ZHIPU_API_KEY}",
        "Content-Type": "application/json"
    }
    data = {
        "model": "glm-4-flash",
        "messages": messages,
        "temperature": 0.1
    }
    resp = requests.post(ZHIPU_API_URL, json=data, headers=headers, timeout=30)
    return resp.json()["choices"][0]["message"]["content"].strip()


# LangChain Prompt 模板（ChatPromptTemplate）
prompt = ChatPromptTemplate.from_messages([
    ("system", "你只能根据下面的知识库内容回答问题。如果知识库里没有相关信息，你必须回答【暂无相关答案】，禁止编造任何信息。"),
    ("user", """知识库内容：
{context}

用户问题：{question}
请严格基于知识库内容回答，不要使用你的背景知识："""),
])

# LCEL 链式调用：prompt → 智谱 LLM → 字符串解析
chain = prompt | RunnableLambda(_call_zhipu) | StrOutputParser()


def generate_answer(question: str, hits: List[Dict]) -> str:
    if not hits:
        return "抱歉，暂无相关答案"

    context = "\n".join([f"知识：{hit['text']}" for hit in hits])

    try:
        return chain.invoke({"context": context, "question": question})
    except Exception as e:
        print("智谱调用失败：", e)
        return f"答案：{hits[0]['text']}"
