package controller

import (
	"ai-rag-system/go-backend/logic"
	"ai-rag-system/go-backend/mysql/models" // 引入你定义的模型
	"ai-rag-system/go-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// chat.go 是后端和 AI 服务之间的桥梁，负责把用户的问题交给 AI 并把结果返回给前端
// ChatHandler 用户聊天接口
// POST /api/v1/chat
func ChatHandler(c *gin.Context) {
	// 1. 获取登录用户（从 JWT 中间件拿）
	userID, exists := c.Get(jwt.CtxtUserIDKey)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. ✅ 直接使用你定义的 models.Chat 接收参数
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, "参数错误："+err.Error())
		return
	}

	// 3. 调用 AI 服务
	answer, err := logic.ChatAI(
		userID.(int64),
		chat.SessionID,
		chat.Question,
	)
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "AI服务异常："+err.Error())
		return
	}

	// 4. 返回前端
	ResponseSuccess(c, gin.H{
		"session_id": chat.SessionID,
		"answer":     answer,
	})
}
