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

type Task struct {
	gorm.Model
	Name      string //本周主要工作
	State     string //任务状态
	Phone     string
	Email     string
	Address   string
	Content   string    `json:"content" gorm:"default:''"`   //工作详情
	Done      bool      `json:"done"`                        //完成情况
	Uploader  string    `json:"uploader"`                    //负责人
	Assistant string    `json:"assistant"`                   // 新增：辅助人
	StartTime time.Time `json:"start_time"`                  // 新增：任务开始时间
	EndTime   time.Time `json:"end_time"`                    // 新增：任务结束时间
	TaskType  string    `json:"task_type" gorm:"default:''"` // 新增：任务类型(巡检/维修等)
	Priority  int       `json:"priority"`                    // 新增：优先级(1-5)
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

// 更新任务
func (u UserController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(400, gin.H{
			"error": "更新失败",
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
	task.Content = updatedTask.Content
	task.Done = updatedTask.Done
	dao.Db.Save(&task)
	c.JSON(200, gin.H{
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
