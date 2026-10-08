package mysql

import (
	"ai-rag-system/go-backend/mysql/models"
	"time"

	"go.uber.org/zap"
)

// CreateKnowledge 创建知识库
func CreateKnowledge(kn *models.Knowledge) error {
	now := time.Now().Unix()
	kn.CreatedAt = now
	kn.UpdatedAt = now

	sqlStr := `INSERT INTO knowledge (
		title, 
		file_name, 
		path, 
		content, 
		created_at, 
		updated_at, 
		source
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(sqlStr,
		kn.Title,
		kn.FileName,
		kn.Path,
		kn.Content,
		kn.CreatedAt,
		kn.UpdatedAt,
		kn.Source,
	)
	return err
}

// GetKnowledgeList 获取列表
func GetKnowledgeList() ([]models.Knowledge, error) {
	var list []models.Knowledge
	sqlStr := `
	SELECT 
		id, 
		title, 
		file_name, 
		path, 
		content, 
		created_at, 
		updated_at, 
		source
	FROM knowledge 
	ORDER BY id DESC`

	zap.L().Info("执行SQL", zap.String("sql", sqlStr))
	err := db.Select(&list, sqlStr)
	return list, err
}

// DeleteKnowledge 删除
func DeleteKnowledge(id int64) error {
	sqlStr := `DELETE FROM knowledge WHERE id=?`
	_, err := db.Exec(sqlStr, id)
	return err
}

// SyncKnowledgeAI 同步AI用
func SyncKnowledgeAI(id int64) (models.Knowledge, error) {
	var kn models.Knowledge
	sqlStr := `
	SELECT 
		id, 
		title, 
		file_name, 
		path, 
		content, 
		source, 
		created_at,
		updated_at
	FROM knowledge 
	WHERE id=?`

	err := db.Get(&kn, sqlStr, id)
	return kn, err
}

// SaveKnowledge 保存并返回ID
func SaveKnowledge(title, fileName, path, content, source string) (int64, error) {
	now := time.Now().Unix()

	query := `
	INSERT INTO knowledge (
		title, 
		file_name, 
		path, 
		content, 
		source, 
		created_at, 
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query,
		title,
		fileName,
		path,
		content,
		source,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetKnowledgeByID 根据ID查询
func GetKnowledgeByID(id int64) (*models.Knowledge, error) {
	var kn models.Knowledge
	err := db.Get(&kn,
		`SELECT 
			id, 
			title, 
			file_name, 
			path, 
			content, 
			source, 
			created_at,
			updated_at
		FROM knowledge 
		WHERE id = ?`,
		id,
	)
	return &kn, err
}
