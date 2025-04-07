package models

import (
	"demo/dao"
	"time"
)

type DailyTask struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Done      bool      `json:"done"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Priority  int       `json:"priority"`
}

func GetDailyTasks(date time.Time) ([]DailyTask, error) {
	var tasks []DailyTask
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	err := dao.Db.Model(&Task{}).
		Where("start_time BETWEEN ? AND ?", start, end).
		Order("priority DESC").
		Find(&tasks).Error

	return tasks, err
}
