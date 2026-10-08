package models

//chat 问答模型
type Chat struct {
	SessionID string `json:"session_id" binding:"required"` // 会话ID
	Question  string `json:"question" binding:"required"`
	Answer    string `json:"answer"`
}
