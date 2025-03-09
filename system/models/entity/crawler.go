// entity/crawler.go
package entity

import (
	"github.com/google/uuid"
)

type Crawler struct {
	BaseModel // 继承基础模型（包含 ID, CreatedAt, UpdatedAt）

	// 基础配置
	Name        string   `json:"name" redis:"name" gorm:"type:varchar(100);uniqueIndex"` // 爬虫名称（唯一索引）
	Description string   `json:"description" redis:"description" gorm:"type:text"`       // 描述
	Status      string   `json:"status" redis:"status" gorm:"type:varchar(20);index"`    // 状态（running/stopped）
	Config      Config   `json:"config" redis:"config" gorm:"type:json"`                 // 爬虫配置（JSON 存储）
	Schedule    Schedule `json:"schedule" redis:"schedule" gorm:"type:json"`             // 调度规则（JSON 存储）

	// 关联信息
	TaskIDs []string `json:"task_ids" redis:"task_ids" gorm:"type:json"`      // 关联的任务 ID 列表
	UserID  string   `json:"-" redis:"user_id" gorm:"type:varchar(36);index"` // 所属用户 ID（不返回前端）
}

// Config 爬虫配置细节
type Config struct {
	StartURLs   []string `json:"start_urls"`      // 初始 URL 列表
	Concurrency int      `json:"concurrency"`     // 并发数
	Proxy       string   `json:"proxy,omitempty"` // 代理地址（可选）
}

// Schedule 调度规则
type Schedule struct {
	Cron       string `json:"cron"`        // Cron 表达式
	RetryTimes int    `json:"retry_times"` // 失败重试次数
}

// 生成 UUID（与 Task 结构保持一致）
func (c *Crawler) BeforeCreate() error {
	c.ID = uuid.New().String()
	return nil
}
