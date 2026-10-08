package controller

import (
	"ai-rag-system/go-backend/logic"
	"ai-rag-system/go-backend/mysql/models"
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// KnowledegHandler 知识库统一入口
func KnowledegHandler(c *gin.Context) {
	action := c.Query("action")

	switch action {
	case "upload":
		UploadFileHandler(c)
	case "save":
		SaveKnowledgeHandler(c)
	case "list":
		GetKnowledgeListHandler(c)
	case "delete":
		DeleteKnowledgeHandler(c)
	default:
		ResponseErrorWithMsg(c, CodeInvalidParam, "无效的action")
	}
}

// 上传文件接口

func UploadFileHandler(c *gin.Context) {
	// 1. 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 创建目录
	uploadDir := "uploaded_files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		zap.L().Error("创建目录失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 处理文件名
	fileName := file.Filename
	ext := filepath.Ext(fileName)
	for strings.HasSuffix(fileName, ext+ext) {
		fileName = strings.TrimSuffix(fileName, ext)
	}

	// 4. 保存文件
	savePath := filepath.Join(uploadDir, fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		zap.L().Error("保存文件失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 5. 读取内容
	contentBytes, err := os.ReadFile(savePath)
	if err != nil {
		zap.L().Error("读取文件失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 6. 组装结构体
	kn := &models.Knowledge{
		Title:    fileName,
		FileName: fileName,
		Path:     savePath,
		Content: sql.NullString{
			String: string(contentBytes),
			Valid:  true,
		},
		Source: "upload",
	}

	// 7. 交给 logic 保存（自动带时间戳 + 自动同步AI ✅）
	if err := logic.SaveKnowledge(kn); err != nil {
		zap.L().Error("保存失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, "上传成功")
}

// 手动保存知识库
func SaveKnowledgeHandler(c *gin.Context) {
	p := new(models.Knowledge)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SaveKnowledge 参数错误", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	if p.Path != "" {
		contentBytes, err := os.ReadFile(p.Path)
		if err == nil {
			p.Content = sql.NullString{
				String: string(contentBytes),
				Valid:  true,
			}
		} else {
			zap.L().Error("读取文件内容失败", zap.String("path", p.Path), zap.Error(err))
		}
	}

	if err := logic.SaveKnowledge(p); err != nil {
		zap.L().Error("保存失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, "保存成功")
}

// 获取列表
func GetKnowledgeListHandler(c *gin.Context) {
	list, err := logic.GetKnowledgeList()
	if err != nil {
		zap.L().Error("获取列表失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccess(c, list)
}

// 删除
func DeleteKnowledgeHandler(c *gin.Context) {
	idStr := c.Query("id")
	if err := logic.DeleteKnowledge(idStr); err != nil {
		zap.L().Error("删除失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, "删除成功")
}
