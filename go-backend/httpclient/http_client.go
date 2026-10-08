package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 全局HTTP客户端（超时+复用）
var (
	AiClient = &http.Client{
		Timeout: 60 * time.Second, // 大文件同步可以长一点
	}
	AiBaseURL string
)

// 初始化：从 config.yaml 读取 python.url
func InitAiHttpClient() {
	AiBaseURL = viper.GetString("python.url")
	if AiBaseURL == "" {
		zap.L().Error("配置文件中 python.url 未配置")
	}
	zap.L().Info("AI HTTP客户端初始化成功", zap.String("ai_url", AiBaseURL))
}

// go连接python的统一HTTP客户端，
// PostJSON
// 统一发送 JSON 格式的 POST 请求给 Python AI 服务
// path: 请求路径，例如 /chat 或 /sync
// reqData: Go 结构体或 map，会自动序列化为 JSON
// 返回：Python 返回的原始数据 []byte 和错误信息
func PostJSON(path string, reqData interface{}) ([]byte, error) {
	// 1. 将 Go 结构体/map 序列化为 JSON 字节流
	// 因为 HTTP 传输只能识别字节流，不能直接传结构体
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err // 序列化失败直接返回错误
	}

	// 2. 拼接完整请求地址
	// AiBaseURL = http://127.0.0.1:8000 ✅ 正确
	// path      = /chat                ✅ 正确
	// fullURL   = http://127.0.0.1:8000/chat ✅ 正确
	fullURL := AiBaseURL + path

	// 3. 创建 POST 请求
	// 参数：请求方法、完整URL、请求体（JSON）
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err // 创建请求失败直接返回
	}

	// 4. 设置请求头
	// 必须告诉 Python：我发送的是 JSON 格式数据
	req.Header.Set("Content-Type", "application/json")

	// 5. 发送请求到 Python AI 服务
	// AiClient 是全局 HTTP 客户端（带超时、复用连接）
	resp, err := AiClient.Do(req)
	if err != nil {
		return nil, err // 请求发送失败（连不上Python）
	}
	// 函数结束时关闭响应体，避免内存泄漏
	defer resp.Body.Close()

	// 6. 读取 Python 返回的全部响应数据
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // 读取失败
	}

	// 7. 判断 HTTP 状态码（200~299 才算成功）
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 记录错误日志
		zap.L().Error("AI服务调用失败",
			zap.Int("status", resp.StatusCode),
			zap.String("url", fullURL),
		)
		// 虽然失败，但仍返回结果，方便调试
		return result, err
	}

	// 8. 成功：返回 Python 给的原始数据
	return result, nil
}
