package snowflake

import (
	"fmt"
	"sync"

	sf "github.com/bwmarrin/snowflake"
)

// 全局唯一的 Node 对象
var node *sf.Node
var once sync.Once

// Init 初始化雪花算法节点
// machineID: 当前机器/服务的唯一ID (0-1023)
// 在分布式系统中，每台部署的服务器必须拥有不同的 machineID，否则会 ID 冲突！
func Init(cfg Config) (err error) {
	// 锁定初始时间，避免 ID 生成随时间回拨导致问题（可选配置）
	// sf.Epoch = 1640995200000 // 例如设置起始时间为 2022-01-01

	once.Do(func() {
		node, err = sf.NewNode(cfg.MachineID)
	})

	if err != nil {
		return fmt.Errorf("init snowflake failed: %v", err)
	}
	return nil
}

func GenID() int64 {
	if node == nil {
		panic("snowflake node not initialized")
	}
	return node.Generate().Int64()
}

func GenStringID() string {
	if node == nil {
		panic("snowflake node not initialized")
	}
	return node.Generate().String()
}
