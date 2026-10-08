package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

// 全局雪花节点变量
var node *snowflake.Node

// Init 初始化雪花算法
// startTime: 起始时间戳 (格式 "2006-01-02")
// machineID: 机器ID (0-1023)
func Init(startTime string, machineID int64) error {
	// ✅ 正确时间格式：2006-01-02
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return fmt.Errorf("时间解析失败: %v", err)
	}

	// 设置雪花算法的 Epoch (起始时间戳)
	snowflake.Epoch = st.UnixMilli()

	// 创建节点
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return fmt.Errorf("创建雪花节点失败: %v", err)
	}

	return nil
}

// GenID 生成雪花算法ID
func GenID() int64 {
	return node.Generate().Int64()
}
