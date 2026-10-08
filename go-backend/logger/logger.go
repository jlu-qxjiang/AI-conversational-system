package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init 初始化日志（从Viper读取配置）
func Init() (err error) {
	// 1. 配置日志写入器（lumberjack 日志切割）日志写到哪里
	writeSyncer := getLogWriter(
		viper.GetString("log.filename"), //日志名称
		viper.GetInt("log.max_size"),    //日志大小
		viper.GetInt("log.max_backups"), //最多备份几个
		viper.GetInt("log.max_age"),     //最多保存几天
	)

	// 2. 配置日志编码器（格式）日志长什么样子
	encoder := getEncoder()

	// 3. 读取并解析日志级别
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return err
	}

	// 4. 创建核心：文件 + 终端 同时输出
	fileCore := zapcore.NewCore(encoder, writeSyncer, l)                   //打印到日志文件
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), l) //打印到终端

	// 合并两个输出
	core := zapcore.NewTee(fileCore, consoleCore)

	// 5. 初始化全局日志（添加调用者信息）
	lg := zap.New(core, zap.AddCaller())
	// 替换zap全局日志，方便全局使用zap.L()
	zap.ReplaceGlobals(lg)

	return nil
}

// getEncoder 配置日志编码格式（日志长什么样子）
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 配置日志写入（lumberjack 日志切割）
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志，将其集成到zap日志系统
// 请求进来 → 记录开始时间 → 执行接口 → 计算耗时 → 打印日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 记录请求开始时间、请求路径、查询参数
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 2. 执行后续中间件/路由处理 放行！去执行后面的中间件、接口业务逻辑
		c.Next()

		// 3. 计算请求耗时，记录请求详情
		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),                                 // 状态码
			zap.String("method", c.Request.Method),                               // 请求方法 GET/POST
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // URL ?后面的参数
			zap.String("ip", c.ClientIP()),                                       // 用户IP
			zap.String("user-agent", c.Request.UserAgent()),                      // 浏览器/客户端信息
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 错误信息
			zap.Duration("cost", cost),                                           // 请求耗时
		)
	}
}

// GinRecovery 捕获panic，打印堆栈日志，避免服务崩溃
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  defer 捕获 panic
		defer func() {
			if err := recover(); err != nil {
				// 判断是否为连接异常（不需要打印堆栈）
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 读取请求详情
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 连接断开，不需要记录堆栈
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Abort()
					return
				}

				// 真正的 panic，打印堆栈信息
				zap.L().Error("panic recovered",
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回 500 给前端
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}
		}()

		// 继续执行后面的中间件
		c.Next()
	}
}
