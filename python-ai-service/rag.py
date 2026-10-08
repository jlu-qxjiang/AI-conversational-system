import os
import json
import time
import numpy as np
import faiss
import requests
from typing import List, Dict, Any

# 智谱 Embedding API
ZHIPU_API_KEY = os.environ.get("ZHIPU_API_KEY", "")
ZHIPU_EMBED_URL = "https://open.bigmodel.cn/api/paas/v4/embeddings"
EMBED_MODEL = "embedding-2"
EMBED_DIM = 1024  # embedding-2 输出 1024 维

INDEX_FILE = "faiss_index.bin"
META_FILE = "docs_meta.json"


class FAISSRAG:
    def __init__(self):
        self.index = faiss.IndexFlatIP(EMBED_DIM)  # 内积 = 余弦相似度（向量已归一化）
        self.docs: List[Dict[str, Any]] = []
        self.load()

    # ======================
    # 调用智谱 Embedding API 生成向量
    # ======================
    def _embed_batch(self, texts: List[str]) -> np.ndarray:
        headers = {
            "Authorization": f"Bearer {ZHIPU_API_KEY}",
            "Content-Type": "application/json"
        }
        data = {
            "model": EMBED_MODEL,
            "input": texts
        }
        # 智谱 embedding API 有并发限制，失败时重试
        for attempt in range(3):
            try:
                resp = requests.post(ZHIPU_EMBED_URL, json=data, headers=headers, timeout=30)
                result = resp.json()
                vecs = []
                for item in result["data"]:
                    v = np.array(item["embedding"], dtype=np.float32)
                    # 归一化以便内积等价余弦相似度
                    norm = np.linalg.norm(v)
                    if norm > 0:
                        v = v / norm
                    vecs.append(v)
                return np.array(vecs)
            except Exception as e:
                if attempt < 2:
                    time.sleep(1)
                else:
                    raise e

    def embed_text(self, text: str) -> np.ndarray:
        return self._embed_batch([text])[0]

    # ======================
    # 添加文档到 FAISS 索引
    # ======================
    def add_documents(self, docs: List[Dict[str, Any]]):
        texts_to_embed = []
        metas_to_add = []
        for doc in docs:
            content = doc.get("content", "")
            if not content:
                continue
            texts_to_embed.append(content)
            metas_to_add.append({
                "id": doc.get("id"),
                "title": doc.get("title"),
                "source": doc.get("source"),
                "text": content
            })

        if not texts_to_embed:
            return

        vecs = self._embed_batch(texts_to_embed)
        self.index.add(vecs)
        self.docs.extend(metas_to_add)
        self.save()

    # ======================
    # 语义检索
    # ======================
    def search(self, query: str, top_k: int = 1, min_score: float = 0.3) -> List[Dict[str, Any]]:
        if self.index.ntotal == 0:
            return []
        q_vec = self.embed_text(query)
        scores, indices = self.index.search(np.array([q_vec]), top_k)
        results = []
        for score, idx in zip(scores[0], indices[0]):
            if 0 <= idx < len(self.docs) and score >= min_score:
                results.append({**self.docs[idx], "_score": float(score)})
        return results

    # ======================
    # 按 ID 删除文档
    # ======================
    def delete_by_id(self, doc_id: str):
        changed = False
        for i in range(len(self.docs) - 1, -1, -1):
            if str(self.docs[i].get("id", "")) == doc_id:
                del self.docs[i]
                changed = True
                break
        if changed:
            self._rebuild_index()
            self.save()
            return True
        return False

    # ======================
    # 全量重建 FAISS 索引
    # ======================
    def _rebuild_index(self):
        self.index = faiss.IndexFlatIP(EMBED_DIM)
        if not self.docs:
            return
        texts = [d["text"] for d in self.docs]
        # 分批处理避免一次性传太多文本
        batch_size = 16
        for i in range(0, len(texts), batch_size):
            batch = texts[i:i + batch_size]
            vecs = self._embed_batch(batch)
            self.index.add(vecs)

    # ======================
    # 持久化：FAISS 索引 + 元数据 JSON
    # ======================
    def save(self):
        faiss.write_index(self.index, INDEX_FILE)
        with open(META_FILE, "w", encoding="utf-8") as f:
            json.dump(self.docs, f, ensure_ascii=False, indent=2)

    def load(self):
        if os.path.exists(INDEX_FILE) and os.path.exists(META_FILE):
            self.index = faiss.read_index(INDEX_FILE)
            with open(META_FILE, "r", encoding="utf-8") as f:
                self.docs = json.load(f)
