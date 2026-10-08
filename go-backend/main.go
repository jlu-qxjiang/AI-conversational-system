package main

import (
	"ai-rag-system/go-backend/controller"
	"ai-rag-system/go-backend/httpclient"
	"ai-rag-system/go-backend/logger"
	"ai-rag-system/go-backend/mysql"
	"ai-rag-system/go-backend/pkg/snowflake"
	"ai-rag-system/go-backend/redis"
	"ai-rag-system/go-backend/routes"
	"ai-rag-system/go-backend/settings"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//初始化配置
	err := settings.Init()
	if err != nil {
		fmt.Println("Failed to initialize settings:", err)
		return
	}
	//初始化日志
	err = logger.Init()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	zap.L().Debug("logger initialized")
	//mysql初始化
	err = mysql.Init()
	if err != nil {
		fmt.Println("Failed to initialize mysql:", err)
		return
	}
	zap.L().Debug("mysql initialized")
	defer mysql.Close()
	zap.L().Debug("mysql closed")
	//redis初始化
	err = redis.Init()
	if err != nil {
		fmt.Println("Failed to initialize redis:", err)
		return
	}
	zap.L().Debug("redis initialized")
	defer redis.Close()
	zap.L().Debug("redis closed")
	// 5. 初始化雪花算法（从配置文件读取）
	if err := snowflake.Init(
		settings.Config.Snowflake.StartTime,
		settings.Config.Snowflake.MachineID,
	); err != nil {
		zap.L().Fatal("init snowflake failed", zap.Error(err))
		return
	}
	zap.L().Debug("snowflake initialized")
	//初始化gin框架内置的校验器使用翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init InitTrans failed:%v\n", err)
	}
	//6. 初始化HTTP客户端发出连接Python的请求
	httpclient.InitAiHttpClient()

	zap.L().Info("=== 开始测试 Python 连接 ===")
	resp, err := httpclient.PostJSON("/test", map[string]string{
		"from": "go_backend",
		"test": "connection",
	})
	if err != nil {
		zap.L().Error("❌ Python 连接失败", zap.Error(err))
	} else {
		var result map[string]interface{}
		json.Unmarshal(resp, &result)
		zap.L().Info("✅ Python 连接成功", zap.Any("response", result))
	}
	// 7. 注册路由
	r := routes.SetupRouter()

	// 8. 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
