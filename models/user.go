package models

import (
	"demo/dao"
	"gorm.io/gorm"
	"time"
)

// dao.go
type Task struct {
	gorm.Model
	Name      string    `json:"name" gorm:"not null"`
	State     string    `json:"state" gorm:"default:'pending'"`
	Phone     string    `json:"phone" gorm:"default:''"`
	Email     string    `json:"email" gorm:"default:''"`
	Address   string    `json:"address" gorm:"default:''"`
	Content   string    `json:"content" `
	Done      bool      `json:"done" gorm:"default:false"`
	Uploader  string    `json:"uploader" gorm:"not null"`
	Assistant string    `json:"assistant" gorm:"default:''"`
	StartTime time.Time `json:"start_time" gorm:"default:CURRENT_TIMESTAMP"`
	EndTime   time.Time `json:"end_time" gorm:"default:(CURRENT_TIMESTAMP + INTERVAL 1 DAY)"`
	TaskType  string    `json:"task_type" gorm:"default:''"`
	Priority  int       `json:"priority" gorm:"default:3"`
}

func (Task) TableName() string {
	return "user"
}

func GetUserTest(id int) (Task, error) {
	var user Task
	err := dao.Db.Where("id = ?", id).First(&user).Error
	return user, err
}
