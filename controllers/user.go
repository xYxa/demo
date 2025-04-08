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

// 增强GenerateWeeklyReport函数，支持HTML格式返回
func (u UserController) GenerateWeeklyReport(c *gin.Context) {
	// 获取查询参数
	format := c.DefaultQuery("format", "json") // 默认为json格式

	// 获取最近一周的任务数据
	var tasks []Task
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	if err := dao.Db.Where("created_at >= ? OR updated_at >= ?", oneWeekAgo, oneWeekAgo).Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取周报数据失败"})
		return
	}

	// 按格式返回
	if format == "html" {
		// 生成HTML格式周报
		htmlReport := generateHTMLReport(tasks)
		c.Header("Content-Type", "text/html")
		c.String(200, htmlReport)
	} else {
		// 默认返回JSON格式
		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"report_title":    fmt.Sprintf("10.运维服务周报(%d年第%d周)", time.Now().Year(), getWeekOfYear()),
				"fill_date":       time.Now().Format("2006年01月02日"),
				"filler":          "系统自动生成",
				"completed_tasks": filterTasks(tasks, true),
				"pending_tasks":   filterTasks(tasks, false),
				"statistics": gin.H{
					"total":       len(tasks),
					"completed":   countCompleted(tasks),
					"in_progress": len(tasks) - countCompleted(tasks),
					"completion":  fmt.Sprintf("%.1f%%", float64(countCompleted(tasks))/float64(len(tasks))*100),
				},
			},
		})
	}
}

// 生成HTML格式周报
// generateHTMLReport generates an Excel-compatible HTML report from tasks
func generateHTMLReport(tasks []Task) string {
	now := time.Now()
	weekNum := getWeekOfYear()
	completedTasks := filterTasks(tasks, true)
	pendingTasks := filterTasks(tasks, false)
	stats := gin.H{
		"total":       len(tasks),
		"completed":   len(completedTasks),
		"in_progress": len(pendingTasks),
		"completion":  fmt.Sprintf("%.1f%%", float64(len(completedTasks))/float64(len(tasks))*100),
	}

	// Generate the HTML template
	html := fmt.Sprintf(`
<html xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="ProgId" content="Excel.Sheet"/>
<meta name="Generator" content="Go Weekly Report"/>
<title>运维周报</title>
<style>
.xl67 { text-align:center; vertical-align:middle; color:windowtext; font-size:18.0pt; font-weight:700; font-family:方正仿宋_GBK; }
.xl68 { text-align:left; vertical-align:middle; color:windowtext; font-size:9.0pt; font-weight:700; border-bottom:.5pt solid windowtext; }
.xl69 { text-align:center; vertical-align:middle; background:#BBE4D6; color:windowtext; font-size:14.0pt; font-weight:700; border:.5pt solid windowtext; }
.xl71 { text-align:center; vertical-align:middle; background:#F6F3F6; color:windowtext; font-size:9.0pt; font-weight:700; border:.5pt solid windowtext; }
.xl75 { text-align:center; vertical-align:middle; background:#F6F3F6; color:windowtext; font-size:9.0pt; border:.5pt solid windowtext; }
.xl76 { text-align:center; vertical-align:middle; color:windowtext; font-size:9.0pt; border:.5pt solid windowtext; }
.xl79 { text-align:center; vertical-align:middle; color:windowtext; font-size:9.0pt; border:.5pt solid windowtext; }
.xl83 { text-align:center; vertical-align:middle; background:#BBE4D6; color:windowtext; font-size:14.0pt; font-weight:700; border:.5pt solid windowtext; }
.xl92 { text-align:center; vertical-align:middle; background:#FFFFFF; color:windowtext; font-size:9.0pt; border:.5pt solid windowtext; }
.xl100 { vertical-align:middle; color:windowtext; font-size:9.0pt; border:.5pt solid windowtext; }
.completed { color:#52c41a; font-weight:bold; }
.pending { color:#faad14; }
.overdue { color:#f5222d; }
</style>
</head>
<body link="blue" vlink="purple">
<table width="1391.25" border="0" cellpadding="0" cellspacing="0" style='width:834.75pt;border-collapse:collapse;table-layout:fixed;'>
<col width="41.67"/>
<col width="22.08"/>
<col width="91.58"/>
<col width="270"/>
<col width="59.67"/>
<col width="116.17"/>
<col width="71.17"/>
<col width="41.67"/>
<col width="26.17"/>
<col width="91.58"/>
<col width="72.83"/>
<col width="89.17"/>
<col width="113.67"/>
<col width="59.67"/>
<col width="97.33"/>
<col width="126.83"/>

<!-- 报表标题 -->
<tr height="130">
<td class="xl67" height="130" colspan="16">运维服务周报<br/>(10.湖南-安化运维周报%d年第%d周)<br/></td>
</tr>

<!-- 填写信息 -->
<tr height="36.50">
<td class="xl68" height="36.50" colspan="16">填写日期：%s&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 报告填写人：余湘</td>
</tr>

<!-- 统计信息 -->
<tr height="40">
<td class="xl69" colspan="4">任务统计</td>
<td class="xl69" colspan="4">总任务数: %d</td>
<td class="xl69" colspan="4">已完成: %d</td>
<td class="xl69" colspan="4">完成率: %s</td>
</tr>

<!-- 本周工作总结和下周工作计划标题 -->
<tr height="33.58">
<td class="xl69" height="33.58" colspan="7">本 周 工 作 总 结</td>
<td class="xl69" colspan="9">下 周 工 作 计 划</td>
</tr>

<!-- 表格列标题 -->
<tr height="40">
<td class="xl71">编号</td>
<td class="xl71" colspan="3">本周主要工作</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
<td class="xl71">完成情况</td>
<td class="xl71">编号</td>
<td class="xl71" colspan="4">下周工作主要内容</td>
<td class="xl71">计划完成时间</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
<td class="xl71">状态</td>
</tr>
`, now.Year(), weekNum, now.Format("2006年01月02日"), stats["total"], stats["completed"], stats["completion"])

	// 添加已完成任务
	for i, task := range completedTasks {
		html += fmt.Sprintf(`
<tr height="41.67">
<td class="xl75">%d</td>
<td class="xl76" colspan="3">%s</td>
<td class="xl79">%s</td>
<td class="xl79">%s</td>
<td class="xl79 completed">已完成</td>
<td class="xl75">%d</td>
<td class="xl76" colspan="4">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
</tr>
`, i+1, task.Name, task.Uploader, task.Assistant, i+1,
			task.Content, task.EndTime.Format("2006-01-02"),
			task.Uploader, task.Assistant, getTaskStatus(task))
	}

	// 添加未完成任务
	for i, task := range pendingTasks {
		statusClass := "pending"
		if time.Now().After(task.EndTime) {
			statusClass = "overdue"
		}

		html += fmt.Sprintf(`
<tr height="41.67">
<td class="xl75">%d</td>
<td class="xl76" colspan="3">%s</td>
<td class="xl79">%s</td>
<td class="xl79">%s</td>
<td class="xl79 %s">%s</td>
<td class="xl75">%d</td>
<td class="xl76" colspan="4">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
<td class="xl92">%s</td>
</tr>
`, len(completedTasks)+i+1, task.Name, task.Uploader, task.Assistant,
			statusClass, getTaskStatus(task), len(completedTasks)+i+1,
			task.Content, task.EndTime.Format("2006-01-02"),
			task.Uploader, task.Assistant, getTaskStatus(task))
	}

	// 添加问题记录部分
	html += `
<tr height="33.58">
<td class="xl83" colspan="16">运维工作遇到的主要问题</td>
</tr>

<tr height="40">
<td class="xl71">编号</td>
<td class="xl71" colspan="11">问题描述</td>
<td class="xl71">是否解决</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
</tr>

<tr height="43.33">
<td class="xl75">1</td>
<td class="xl76" colspan="11">无重大问题</td>
<td class="xl100">已解决</td>
<td class="xl92">余湘</td>
<td class="xl92">-</td>
</tr>

<tr height="29">
<td class="xl83" colspan="16">意见及建议</td>
</tr>

<tr height="24">
<td class="xl71">编号</td>
<td class="xl71" colspan="14">内容</td>
<td class="xl71">提出人</td>
</tr>

<tr height="24">
<td class="xl75">1</td>
<td class="xl76" colspan="14"> &nbsp;&nbsp;</td>
<td class="xl92"> &nbsp;&nbsp;</td>
</tr>
</table>
</body>
</html>
`

	return html
}

// 辅助函数：过滤任务
func filterTasks(tasks []Task, completed bool) []Task {
	var result []Task
	for _, task := range tasks {
		if task.Done == completed {
			result = append(result, task)
		}
	}
	return result
}

// 辅助函数：统计已完成任务
func countCompleted(tasks []Task) int {
	count := 0
	for _, task := range tasks {
		if task.Done {
			count++
		}
	}
	return count
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
