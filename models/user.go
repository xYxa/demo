package models

import (
	"demo/dao"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name      string //本周主要工作
	State     string //任务状态
	Phone     string
	Email     string
	Address   string
	Content   string    `json:"content"`    //工作详情
	Done      bool      `json:"done"`       //完成情况
	Uploader  string    `json:"uploader"`   //负责人
	Assistant string    `json:"assistant"`  // 新增：辅助人
	StartTime time.Time `json:"start_time"` // 新增：任务开始时间
	EndTime   time.Time `json:"end_time"`   // 新增：任务结束时间
	TaskType  string    `json:"task_type"`  // 新增：任务类型(巡检/维修等)
	Priority  int       `json:"priority"`   // 新增：优先级(1-5)
}

func (Task) TableName() string {
	return "user"
}

func GetUserTest(id int) (Task, error) {
	var user Task
	err := dao.Db.Where("id = ?", id).First(&user).Error
	return user, err
}
