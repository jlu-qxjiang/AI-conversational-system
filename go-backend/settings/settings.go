package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 对配置文件进行结构体的说明，方便后续调用
type AppConfig struct {
	App struct {
		Name    string `mapstructure:"name"`
		Mode    string `mapstructure:"mode"`
		Port    int    `mapstructure:"port"` //端口号
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`
	MySQL struct {
		Host         string `mapstructure:"host"`           //主机地址	默认值为localhost
		Port         int    `mapstructure:"port"`           //端口号
		User         string `mapstructure:"user"`           //用户名
		Password     string `mapstructure:"password"`       //密码
		DBName       string `mapstructure:"dbname"`         //数据库名
		MaxOpenConns int    `mapstructure:"max_open_conns"` //最大打开连接数	默认值为0，表示不限制连接数
		MaxIdleConns int    `mapstructure:"max_idle_conns"` //最大空闲连接数	默认值为0，表示不限制连接数
	} `mapstructure:"mysql"`
	Redis struct {
		Host     string `mapstructure:"host"`      //主机地址	默认值为localhost
		Port     int    `mapstructure:"port"`      //端口号
		Password string `mapstructure:"password"`  //密码
		DB       int    `mapstructure:"db"`        //数据库编号
		PoolSize int    `mapstructure:"pool_size"` //连接池大小	默认值为0，表示不限制连接数
	} `mapstructure:"redis"`
	Log struct {
		Level      string `mapstructure:"level"`       //日志级别
		Filename   string `mapstructure:"filename"`    //日志文件名
		MaxSize    int    `mapstructure:"max_size"`    //日志文件最大大小（MB）
		MaxAge     int    `mapstructure:"max_age"`     //日志文件最大保存天数
		MaxBackups int    `mapstructure:"max_backups"` //最大备份数
	} `mapstructure:"log"`
	Snowflake struct {
		StartTime string `mapstructure:"start_time"`
		MachineID int64  `mapstructure:"machine_id"`
	} `mapstructure:"snowflake"`
}

var Config = new(AppConfig)

func Init() (err error) {
	//1.设置配置文件路径
	viper.SetConfigFile("config/config.yaml")
	//2.加载配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %w \n", err)
	}
	//3.绑定结构体
	if err := viper.Unmarshal(Config); err != nil {
		return fmt.Errorf("Unmarshal config error: %w \n", err)
	}
	// 4. 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
	return nil
}
