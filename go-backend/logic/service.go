package logic

import (
	"ai-rag-system/go-backend/httpclient"
	"ai-rag-system/go-backend/mysql/models"
	"encoding/json"

	"go.uber.org/zap"
)

// ChatAI 统一使用 httpclient 调用 Python AI
func ChatAI(userID int64, sessionID, question string) (string, error) {
	// 构造请求体
	req := models.Chat{

		SessionID: sessionID,
		Question:  question,
	}

	// ✅ 真正调用 Python：使用你写的统一 HTTP 客户端
	respBytes, err := httpclient.PostJSON("/chat", req)
	if err != nil {
		zap.L().Error("调用AI聊天失败", zap.Error(err))
		return "", err
	}

	// 解析 Python 返回结果
	var result struct {
		Answer string `json:"answer"`
	}
	// 把返回的 json 解析成 go 结构体
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return "", err
	}

	return result.Answer, nil
}
