package routes

import (
	"ai-rag-system/go-backend/controller"
	"ai-rag-system/go-backend/logger"
	"ai-rag-system/go-backend/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery())
	r.Use(middlewares.RateLimitMiddleware(1*time.Second, 1000))
	//  跨域配置
	r.Use(middlewares.Cors())
	r.Static("/page", "C:\\Users\\姜泉行\\Desktop\\ai-rag-progect\\page")

	v1 := r.Group("/ai/v1")

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// 用户注册
	v1.POST("/signup", controller.SignUpHandler)
	// 用户登录
	v1.POST("/login", controller.LoginHandler)
	//使用token
	v1.Use(middlewares.JWTAuthMiddleware())
	{ //对话
		v1.POST("/chat", controller.ChatHandler)
		// 管理员接口
		v1.POST("/admin", controller.AdminHandler)
		// 知识库接口
		v1.POST("/knowledge", controller.KnowledegHandler) // 上传、删除、保存
		v1.GET("/knowledge", controller.KnowledegHandler)  // 获取列表

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
