package entity

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        string    `json:"id" redis:"id" gorm:"primaryKey;type:varchar(36)"`    // Redis Key 使用此 ID，MySQL 主键
	CreatedAt time.Time `json:"created_at" redis:"created_at" gorm:"autoCreateTime"` // 自动填充时间
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at" gorm:"autoUpdateTime"` // 自动更新
}

type Task struct {
	BaseModel // 内嵌基础模型

	// 基础字段
	CrawlerID  string `json:"spider_id" redis:"spider_id" gorm:"type:varchar(36);index"`       // 关联  表
	Status     string `json:"status" redis:"status" gorm:"type:varchar(32);default:'pending'"` // 任务状态
	NodeID     string `json:"node_id" redis:"node_id" gorm:"type:varchar(36);index"`           // 执行节点 ID
	Cmd        string `json:"cmd" redis:"cmd" gorm:"type:text"`                                // 执行命令
	Param      string `json:"param" redis:"param" gorm:"type:text"`                            // 参数
	Error      string `json:"error" redis:"error" gorm:"type:text"`                            // 错误信息
	Pid        int    `json:"pid" redis:"pid" gorm:"type:int"`                                 // 进程 ID
	ScheduleID string `json:"schedule_id" redis:"schedule_id" gorm:"type:varchar(36)"`         // 调度计划 ID

	Mode     string `json:"mode" redis:"mode" gorm:"type:varchar(32)"`           // 执行模式
	Priority int    `json:"priority" redis:"priority" gorm:"type:int;default:0"` // 优先级

	// 关联字段（Redis 用 JSON 存储，MySQL 用外键或 JSON 字段）
	NodeIDs  []string  `json:"node_ids" redis:"node_ids" gorm:"type:json"`                // 分配的节点 ID 列表
	ParentID string    `json:"parent_id" redis:"parent_id" gorm:"type:varchar(36);index"` // 父任务 ID
	Stat     *TaskStat `json:"stat" redis:"stat" gorm:"type:json"`                        // 统计信息
	SubTasks []Task    `json:"sub_tasks" redis:"sub_tasks" gorm:"-"`                      // 子任务（Redis 存储 JSON，MySQL 不直接存储）
	Crawler  *Crawler  `json:"spider" redis:"spider" gorm:"-"`                            // 关联爬虫信息（非数据库字段）
	UserID   string    `json:"-" redis:"user_id" gorm:"type:varchar(36);index"`           // 用户 ID（不返回给前端）
}

type TaskStat struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int64     `json:"duration"`    // 单位：秒
	MemoryUsed int64     `json:"memory_used"` // 单位：MB
}

// 生成 UUID 作为唯一 ID
func (t *Task) BeforeCreate() error {
	t.ID = uuid.New().String()
	return nil
}
