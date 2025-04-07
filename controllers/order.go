package controllers

import (
	"demo/dao"
	"github.com/gin-gonic/gin"
	"time"
)

type OrderController struct {
}

type Search struct {
	Name string
	Cid  int
}

func (o OrderController) GetList(c *gin.Context) {
	search := &Search{}
	err := c.BindJSON(&search)
	if err == nil {
		ReturnSuccess(c, 0, search.Name, search.Cid, 1)
		return
	}
	ReturnError(c, 4001, gin.H{"err": err})
}

// 新增每日任务查询接口
func (o OrderController) GetDailyTasks(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "日期格式错误，请使用YYYY-MM-DD格式"})
		return
	}

	startTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 0, 1)

	var tasks []Task
	if err := dao.Db.
		Where("start_time >= ? AND end_time <= ?", startTime, endTime).
		Order("priority DESC").
		Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(200, gin.H{
		"date":  dateStr,
		"tasks": tasks,
		"total": len(tasks),
	})
}
