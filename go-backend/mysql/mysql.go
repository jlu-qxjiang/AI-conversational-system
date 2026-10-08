package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 声明全局变量db
var db *sqlx.DB

func Init() (err error) {
	//连接的信息
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"), // ✅ 端口必须是 int
		viper.GetString("mysql.dbname"),
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("mysql connect failed ", zap.Error(err))
		return
	}
	//设置与数据库建立连接的最大数目
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//设置连接池中的最大闲置连接数
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))

	return
}

// 防止db是空指针的时候报错
func Close() {
	if db != nil {
		_ = db.Close()
	}
}
