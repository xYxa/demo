package controllers

import (
	"demo/dao"
	"demo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
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
	Content   string    `json:"content" gorm:"default:'';type:text"`
	Done      bool      `json:"done" gorm:"default:false"`
	Uploader  string    `json:"uploader" gorm:"not null"`
	Assistant string    `json:"assistant" gorm:"default:''"`
	StartTime time.Time `json:"start_time" gorm:"default:CURRENT_TIMESTAMP"`
	EndTime   time.Time `json:"end_time" gorm:"default:(CURRENT_TIMESTAMP + INTERVAL 1 DAY)"`
	TaskType  string    `json:"task_type" gorm:"default:''"`
	Priority  int       `json:"priority" gorm:"default:3"`
}
type UserController struct{}

func (u UserController) GetAdd(c *gin.Context) {
	var data Task
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  "添加失败",
			"data": gin.H{},
			"code": 400,
		})
	} else {
		dao.Db.Create(&data)
		c.JSON(200, gin.H{
			"msg":  "添加成功",
			"data": data,
			"code": 200,
		})
	}
}

func (u UserController) GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	name := c.Param("name")
	id, _ := strconv.Atoi(idStr)
	user, _ := models.GetUserTest(id)
	ReturnSuccess(c, 0, name, user, 1)

}

// 获取所有任务
func (u UserController) GetTasks(c *gin.Context) {
	var tasks []Task
	if err := dao.Db.Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "获取任务列表失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": tasks,
	})
}

// 创建任务
func (u UserController) CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 从 URL 参数中获取上传人信息
	uploader := c.Query("uploader")
	fmt.Println("从 URL 参数中获取的上传人:", uploader) // 打印上传人信息

	if uploader == "" {
		uploader = "匿名用户" // 默认值
	}
	task.Uploader = uploader
	fmt.Println("接收到的任务数据:", task)

	// 将任务保存到数据库
	if err := dao.Db.Create(&task).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "创建任务失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "创建成功",
		"data":    task,
	})

}

// UpdateTask 更新任务（简化版，保持原有路由设计）
func (u UserController) UpdateTask(c *gin.Context) {
	// 1. 获取任务ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "无效的任务ID"})
		return
	}

	// 2. 查询现有任务
	var task Task
	if err := dao.Db.First(&task, id).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": "任务未找到"})
		return
	}

	// 3. 绑定更新数据
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "无效的请求数据"})
		return
	}

	// 4. 保存更新
	if err := dao.Db.Save(&task).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "更新任务失败"})
		return
	}

	// 5. 返回响应（保持原有格式）
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": task,
	})
}

// 删除任务
func (u UserController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // 将字符串 ID 转换为整数
	if err != nil {
		c.JSON(400, gin.H{
			"error": "无效的任务 ID",
		})
		return
	}

	var task Task
	if err := dao.Db.First(&task, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "任务未找到",
		})
		return
	}

	if err := dao.Db.Delete(&task).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "删除任务失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除成功",
	})
}

// GenerateWeeklyReport 生成符合Excel模板格式的运维周报
func (u UserController) GenerateWeeklyReport(c *gin.Context) {
	// 获取最近一周的任务数据
	var tasks []Task
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	if err := dao.Db.Where("created_at >= ?", oneWeekAgo).Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "获取周报数据失败",
		})
		return
	}

	// 准备周报数据
	reportData := gin.H{
		"report_title":   fmt.Sprintf("运维服务周报(10湖南-安化运维周报%d年第%d周)", time.Now().Year(), getWeekOfYear()),
		"fill_date":      time.Now().Format("2006年01月02日"),
		"filler":         "系统自动生成",
		"current_week":   formatWeekTasks(tasks),
		"next_week_plan": []gin.H{},
		"problems":       []gin.H{},
		"next_problems":  []gin.H{},
		"suggestions":    []gin.H{},
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    reportData,
	})
}

// 辅助函数：格式化本周工作任务
func formatWeekTasks(tasks []Task) []gin.H {
	var formattedTasks []gin.H

	for i, task := range tasks {
		formattedTask := gin.H{
			"id":         i + 1,
			"work":       task.Name,
			"leader":     task.Uploader,
			"assistant":  task.Assistant,
			"status":     getTaskStatus(task),
			"start_time": task.StartTime.Format("2006-01-02"),
			"end_time":   task.EndTime.Format("2006-01-02"),
		}
		formattedTasks = append(formattedTasks, formattedTask)
	}

	return formattedTasks
}

// 辅助函数：获取任务状态
func getTaskStatus(task Task) string {
	if task.Done {
		return "已完成"
	}
	if time.Now().After(task.EndTime) {
		return "超期未完成"
	}
	return "进行中"
}

// 辅助函数：获取当前是第几周
func getWeekOfYear() int {
	_, week := time.Now().ISOWeek()
	return week
}
