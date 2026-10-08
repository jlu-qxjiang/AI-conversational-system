package logic

import (
	"ai-rag-system/go-backend/httpclient"
	"ai-rag-system/go-backend/mysql"
	"ai-rag-system/go-backend/mysql/models"
	"strconv"

	"go.uber.org/zap"
)

func GetKnowledgeList() ([]models.Knowledge, error) {
	zap.L().Info("进入logic.GetKnowledgeList")
	list, err := mysql.GetKnowledgeList()
	if err != nil {
		zap.L().Error("mysql查询失败", zap.Error(err))
	} else {
		zap.L().Info("mysql查询成功", zap.Int("条数", len(list)))
	}
	return list, err
}

// 🔥 修复：改用我们新的 SaveKnowledge（能返回 ID）
func SaveKnowledge(kn *models.Knowledge) error {
	// 调用新的保存方法，自动带时间戳，并返回 ID
	id, err := mysql.SaveKnowledge(
		kn.Title,
		kn.FileName,
		kn.Path,
		kn.Content.String,
		kn.Source,
	)
	if err != nil {
		return err
	}

	// 这里才有正确的 ID！
	idStr := strconv.FormatInt(id, 10)
	//向python同步数据
	err = SyncKnowledgeAI(idStr)
	if err != nil {
		zap.L().Error("自动同步AI失败", zap.String("id", idStr), zap.Error(err))
	}

	return nil
}

func DeleteKnowledge(idStr string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	_ = DeleteKnowledgeAI(idStr)
	return mysql.DeleteKnowledge(id)
}

func SyncKnowledgeAI(idStr string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	// 通过id查询完整数据
	kn, err := mysql.GetKnowledgeByID(id)
	if err != nil {
		zap.L().Error("自动同步AI失败：查询失败", zap.String("id", idStr), zap.Error(err))
		return err
	}

	// 准备同步内容
	content := ""
	if kn.Content.Valid {
		content = kn.Content.String
	}

	reqBody := map[string]interface{}{
		"id":         kn.ID,
		"title":      kn.Title,
		"file_name":  kn.FileName,
		"path":       kn.Path,
		"content":    content,
		"source":     kn.Source,
		"created_at": kn.CreatedAt,
	}

	// 调用 Python
	resp, err := httpclient.PostJSON("/sync", reqBody)
	if err != nil {
		zap.L().Error("调用Python同步失败", zap.Error(err))
		return err
	}

	zap.L().Info("Python同步完成", zap.String("resp", string(resp)))
	return nil
}

func DeleteKnowledgeAI(idStr string) error {
	reqBody := map[string]interface{}{
		"id": idStr,
	}
	_, err := httpclient.PostJSON("/delete", reqBody)
	return err
}

func GetKnowledgeByID(id int64) (*models.Knowledge, error) {
	return mysql.GetKnowledgeByID(id)
}
