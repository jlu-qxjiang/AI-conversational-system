package controller

// 导入需要用到的包
import (
	"fmt"           // 用于格式化字符串（生成文件名）
	"io"            // 用于复制文件内容
	"os"            // 用于操作文件、创建文件夹
	"path/filepath" // 用于处理文件路径、文件名
	"time"          // 用于生成时间戳（防止文件名重复）

	"ai-rag-system/go-backend/pkg/jwt" // JWT 工具（拿用户角色）

	"github.com/gin-gonic/gin" // Gin 框架核心
)

// ==============================================
// 函数作用：管理员上传知识库文件
// 请求方式：POST /ai/v1/admin
// 必须传：form-data 格式的 file 文件
// ==============================================
func AdminHandler(c *gin.Context) {

	// ===================== 第一步：校验是不是管理员 =====================
	// 从 Gin 上下文里拿出当前登录用户的角色（role=1是管理员）
	roleValue, exists := c.Get(jwt.CtxtRoleKey)

	// 如果拿不到角色 → 说明没登录 / Token无效
	if !exists {
		ResponseError(c, CodeNeedLogin) // 返回：需要登录
		return
	}

	// 把角色值转成 int 类型
	role, ok := roleValue.(int)

	// 如果转换失败 OR 角色不是 1 → 权限不足
	if !ok || role != 1 {
		ResponseErrorWithMsg(c, CodeInvalidtoken, "管理员权限不足")
		return
	}

	// ===================== 第二步：接收前端上传的文件 =====================
	// 从请求里获取名为 "file" 的文件
	file, header, err := c.Request.FormFile("file")

	// 如果获取文件失败（比如没传文件、字段名不对）
	if err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, "上传文件失败："+err.Error())
		return
	}
	defer file.Close() // 函数结束时自动关闭文件流（防止资源泄漏）

	// ===================== 第三步：创建保存文件的目录 =====================
	// 定义文件夹名称：当前目录下的 uploaded_files
	uploadDir := "./uploaded_files"

	// 创建目录（如果不存在就创建，存在就跳过）
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "创建目录失败："+err.Error())
		return
	}

	// ===================== 第四步：生成不重复的文件名 =====================
	// 文件名规则：时间戳_原文件名
	// 例：1777811284340166000_test.txt
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))

	// 拼接最终保存路径：uploaded_files/xxx.txt
	dstPath := filepath.Join(uploadDir, fileName)

	// ===================== 第五步：创建文件并写入内容 =====================
	// 在服务器上创建一个空文件
	out, err := os.Create(dstPath)
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "保存文件失败："+err.Error())
		return
	}
	defer out.Close() // 结束时关闭文件

	// 把上传的文件内容 复制 到服务器的新文件里
	if _, err := io.Copy(out, file); err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "保存文件失败："+err.Error())
		return
	}

	// ===================== 第六步：返回上传成功结果 =====================
	ResponseSuccess(c, gin.H{
		"file_name": header.Filename, // 原始文件名
		"save_path": dstPath,         // 服务器保存路径
	})
}
