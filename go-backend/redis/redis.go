package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// Init 初始化Redis连接
func Init() (err error) {
	// 1. 创建 Redis 客户端，配置连接参数
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	// 2. 发送 Ping 命令，校验连接是否成功（修复：传入上下文）
	// 生产环境推荐带超时，避免卡死
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //ctx：超时时间、取消信号、请求生命周期
	defer cancel()                                                          // 必须加，防止上下文泄漏
	//向redis发送一个信号检查连接是否成功
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		zap.L().Error("connect redis failed", zap.Error(err))
		return err
	}

	zap.L().Info("connect redis success")
	return nil
}

// GetRdb 暴露全局客户端，供其他包调用
func GetRdb() *redis.Client {
	return rdb
}

// 防止rdb是空指针时报错直接关闭
func Close() {
	_ = rdb.Close()
}
