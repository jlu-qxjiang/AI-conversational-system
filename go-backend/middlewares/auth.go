package middlewares

//实现用户认证
import (
	"ai-rag-system/go-backend/controller"
	"ai-rag-system/go-backend/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		//c → 当前请求的上下文c.Request → HTTP 请求原始对象Header → 请求头→ 拿到 key 为 Authorization 的值
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			//后面的所有环节（中间件、控制器）全部不执行
			c.Abort()
			return
		}
		// 按空格分割
		//前端会传来Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9的token 将其分离成两段
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidtoken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidtoken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(jwt.CtxtUserIDKey, mc.UserID)
		c.Set(jwt.CtxtUsernameKey, mc.Username)
		c.Set(jwt.CtxtRoleKey, mc.Role)
		c.Next() // 后续的处理函数可以用过c.Get(jwt.CtxtUserIDKey)来获取当前请求的用户信息
	}
}
