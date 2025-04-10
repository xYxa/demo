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
	StartTime time.Time `json:"start_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	EndTime   time.Time `json:"end_time" gorm:"type:timestamp;default:(CURRENT_TIMESTAMP + INTERVAL 1 DAY)"`
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
	if uploader == "" {
		uploader = "匿名用户" // 默认值
	}
	task.Uploader = uploader

	// 设置默认时间（如果前端没有提供）
	if task.StartTime.IsZero() {
		task.StartTime = time.Now()
	}
	if task.EndTime.IsZero() {
		task.EndTime = time.Now().Add(24 * time.Hour) // 默认1天后
	}

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

	// 确保时间不会被覆盖为0值
	if task.StartTime.IsZero() {
		task.StartTime = time.Now()
	}
	if task.EndTime.IsZero() {
		task.EndTime = time.Now().Add(24 * time.Hour)
	}

	// 4. 保存更新
	if err := dao.Db.Save(&task).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "更新任务失败"})
		return
	}

	// 5. 返回响应
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
	format := c.DefaultQuery("format", "json")

	// 获取最近一周的任务数据
	var tasks []Task
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	// 打印查询条件
	fmt.Printf("[DEBUG] 查询条件: created_at >= %v OR updated_at >= %v\n",
		oneWeekAgo.Format("2006-01-02 15:04:05"),
		oneWeekAgo.Format("2006-01-02 15:04:05"))

	if err := dao.Db.Where("created_at >= ? OR updated_at >= ?", oneWeekAgo, oneWeekAgo).Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取周报数据失败"})
		return
	}

	// 打印获取到的任务总数
	fmt.Printf("[DEBUG] 获取到 %d 条任务数据\n", len(tasks))

	// 打印每条任务的详细信息
	for i, task := range tasks {
		fmt.Printf("[DEBUG] 任务%d: ID=%d, Name=%s, State=%s, Start=%v, End=%v, Done=%v\n",
			i+1, task.ID, task.Name, task.State,
			task.StartTime.Format("2006-01-02"),
			task.EndTime.Format("2006-01-02"),
			task.Done)
	}

	// 按格式返回
	if format == "html" {
		htmlReport := generateHTMLReport(tasks)
		c.Header("Content-Type", "text/html")
		c.String(200, htmlReport)
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"report_title":    fmt.Sprintf("10.运维服务周报(%d年第%d周)", time.Now().Year(), getWeekOfYear()),
				"completed_tasks": filterTasks(tasks, true),
				"pending_tasks":   filterTasks(tasks, false),
			},
		})
	}
}

// 新增辅助函数：获取下一个周一的日期
func getNextMonday(now time.Time) time.Time {
	weekday := now.Weekday()
	daysUntilMonday := (8 - int(weekday)) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	return now.AddDate(0, 0, daysUntilMonday)
}

// 新增辅助函数：获取下周一到周五的日期列表
func getNextWeekDates(now time.Time) []time.Time {
	nextMonday := getNextMonday(now)
	dates := make([]time.Time, 5)
	for i := 0; i < 5; i++ {
		dates[i] = nextMonday.AddDate(0, 0, i)
	}
	return dates
}

// 判断任务是否属于本周
func isTaskInWeek(task Task, thisMonday, nextMonday time.Time) bool {
	// 确保时间在同一时区
	loc := thisMonday.Location()
	start := task.StartTime.In(loc)
	end := task.EndTime.In(loc)

	return (task.CreatedAt.After(thisMonday) && task.CreatedAt.Before(nextMonday)) ||
		(task.UpdatedAt.After(thisMonday) && task.UpdatedAt.Before(nextMonday)) ||
		(start.Before(nextMonday) && end.After(thisMonday)) || // 跨周任务
		(start.After(thisMonday) && start.Before(nextMonday)) ||
		(end.After(thisMonday) && end.Before(nextMonday))
}

// 生成HTML格式周报
func generateHTMLReport(tasks []Task) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	weekNum := getWeekOfYear()
	thisMonday := getThisMonday(now)
	nextMonday := getNextMonday(now)

	// 调试输出时间范围
	fmt.Printf("[DEBUG] 本周时间范围（本地时区）: %s 至 %s\n",
		thisMonday.Format("2006-01-02 15:04:05"),
		nextMonday.Format("2006-01-02 15:04:05"))

	// 筛选任务
	var thisWeekTasks []Task
	for _, task := range tasks {
		task.StartTime = task.StartTime.In(loc)
		task.EndTime = task.EndTime.In(loc)
		if isTaskInWeek(task, thisMonday, nextMonday) {
			fmt.Printf("[DEBUG] 任务纳入: ID=%d, Name=%s (Start=%s, End=%s)\n",
				task.ID, task.Name,
				task.StartTime.Format("2006-01-02"),
				task.EndTime.Format("2006-01-02"))
			thisWeekTasks = append(thisWeekTasks, task)
		}
	}

	// 分离本周任务和下周任务
	var currentWeekTasks []Task
	var nextWeekTasks []Task

	for _, task := range thisWeekTasks {
		if task.EndTime.After(nextMonday) {
			nextWeekTasks = append(nextWeekTasks, task)
		} else {
			currentWeekTasks = append(currentWeekTasks, task)
		}
	}

	// 获取下周一到周五的日期（不再用于自动填充）
	_ = getNextWeekDates(now)

	fmt.Printf("本周任务数量: %d\n", len(thisWeekTasks))
	for i, task := range thisWeekTasks {
		fmt.Printf("任务%d: Name=%s, Content=%s, Start=%v, End=%v\n",
			i+1, task.Name, task.Content, task.StartTime, task.EndTime)
	}

	// 开始编辑HTML
	html := `<!DOCTYPE html>
<html xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="ProgId" content="Excel.Sheet"/>
<meta name="Generator" content="WPS Office ET"/>
<title>运维周报</title>
<style>
<!-- 
@page {margin:0.43in 0.70in 0.12in 1.18in; mso-header-margin:0.30in; mso-footer-margin:0.08in; mso-page-orientation:landscape;}
.font30
	{color:#FF0000;
	font-size:8.0pt;
	font-weight:700;
	font-style:normal;
	text-decoration:none;
	font-family:"宋体";
	mso-generic-font-family:auto;
	mso-font-charset:134;}
.xl67 {text-align:center; vertical-align:middle; color:windowtext; font-size:18.0pt; font-weight:700; font-family:方正仿宋_GBK; mso-generic-font-family:script; mso-font-charset:134;}
.xl68 {text-align:left; vertical-align:middle; color:windowtext; font-size:9.0pt; font-weight:700; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border-bottom:.5pt solid windowtext;}
.xl69 {text-align:center; vertical-align:middle; background:#BBE4D6; color:windowtext; font-size:14.0pt; font-weight:700; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl70 {vertical-align:middle; background:#BBE4D6; color:windowtext; font-size:12.0pt; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl71 {text-align:center; vertical-align:middle; background:#F6F3F6; color:windowtext; font-size:9.0pt; font-weight:700; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl72 {text-align:center; vertical-align:middle; background:#F6F3F6; color:windowtext; font-size:9.0pt; font-weight:700; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border-left:.5pt solid windowtext; border-top:.5pt solid windowtext; border-bottom:.5pt solid windowtext;}
.xl75 {text-align:center; vertical-align:middle; background:#F6F3F6; color:windowtext; font-size:9.0pt; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl76 {text-align:center; vertical-align:middle; color:windowtext; font-size:9.0pt; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border-left:.5pt solid windowtext; border-top:.5pt solid windowtext; border-bottom:.5pt solid windowtext;}
.xl79 {text-align:center; vertical-align:middle; color:windowtext; font-size:9.0pt; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl83 {text-align:center; vertical-align:middle; background:#BBE4D6; color:windowtext; font-size:14.0pt; font-weight:700; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border-left:.5pt solid windowtext; border-top:.5pt solid windowtext; border-bottom:.5pt solid windowtext;}
.xl91
	{mso-style-parent:style64;
	mso-number-format:"yyyy\0022年\0022m\0022月\0022d\0022日\0022";
	text-align:center;
	vertical-align:middle;
	color:windowtext;
	font-size:9.0pt;
	mso-font-charset:134;
	border:.5pt solid windowtext;}

.xl92 {text-align:center; vertical-align:middle; background:#FFFFFF; color:windowtext; font-size:9.0pt; font-family:宋体; mso-generic-font-family:auto; mso-font-charset:134; border:.5pt solid windowtext;}
.xl93
	{mso-style-parent:style64;
	text-align:center;
	vertical-align:middle;
	mso-pattern:auto none;
	background:#FFFFFF;
	color:windowtext;
	font-size:9.0pt;
	mso-font-charset:134;
	border:.5pt solid windowtext;}
.xl99
	{mso-style-parent:style65;
	text-align:center;
	vertical-align:middle;
	white-space:normal;
	mso-pattern:auto none;
	background:#F6F3F6;
	color:windowtext;
	font-size:9.0pt;
	font-weight:700;
	mso-font-charset:134;
	border:.5pt solid windowtext;}
.xl100
	{mso-style-parent:style64;
	text-align:center;
	vertical-align:middle;
	white-space:normal;
	color:windowtext;
	font-size:9.0pt;
	mso-font-charset:134;
	border:.5pt solid windowtext;}

-->
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
<td class="xl67" height="130" colspan="16">运维服务周报<br/>(10湖南-安化运维周报%d年第%d周)<br/></td>
</tr>

<!-- 填写信息 -->
<tr height="36.50">
<td class="xl68" height="36.50" colspan="16">填写日期：%s&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 报告填写人：余湘</td>
</tr>

<!-- 本周工作总结和下周工作计划标题 -->
<tr height="33.58">
<td class="xl69" height="33.58" colspan="7">本 周 工 作 总 结</td>
<td class="xl69" colspan="9">下 周 工 作 计 划</td>
</tr>

<!-- 表格列标题 -->
<tr height="40">
<td class="xl71">编号</td>
<td class="xl72" colspan="3">本周主要工作</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
<td class="xl71">完成情况</td>
<td class="xl71">编号</td>
<td class="xl72" colspan="4">下周工作主要内容</td>
<td class="xl71">计划完成时间</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
<td class="xl71">是否需要协助<br/><font class="font30">(下拉菜单选择)</font></td>
</tr>`

	// 填充本周工作总结
	html += `<!-- 表格标题行 -->`
	for i, task := range currentWeekTasks {
		html += fmt.Sprintf(`
    <tr height="41.67">
        <td class="xl75">%d</td>
        <td class="xl76" colspan="3">%s</td>
        <td class="xl79">%s</td>
        <td class="xl79">%s</td>
        <td class="xl79 %s">%s</td>
        <td class="xl75">-</td>  <!-- 下周部分留空 -->
        <td class="xl76" colspan="4">-</td>
        <td class="xl91">-</td>
        <td class="xl79">-</td>
        <td class="xl79">-</td>
        <td class="xl92">-</td>
    </tr>`,
			i+1,
			task.Name,
			task.Uploader,
			task.Assistant,
			getTaskStatusClass(task),
			getTaskStatus(task))
	}

	// 填充下周工作计划
	html += `<!-- 下周计划标题行 -->`
	for i, task := range nextWeekTasks {
		html += fmt.Sprintf(`
    <tr height="41.67">
        <td class="xl75">-</td>
        <td class="xl76" colspan="3">-</td>
        <td class="xl79">-</td>
        <td class="xl79">-</td>
        <td class="xl79 pending">进行中</td>
        <td class="xl75">%d</td>
        <td class="xl76" colspan="4">%s</td>  <!-- 显示任务名称 -->
        <td class="xl91">%s</td>
        <td class="xl79">%s</td>
        <td class="xl79">%s</td>
        <td class="xl92">%s</td>
    </tr>`,
			i+1,
			task.Name, // 确保显示任务名称
			task.EndTime.Format("2006年01月02日"),
			task.Uploader,
			task.Assistant,
			getAssistanceNeeded(task))
	}

	// Add problems section
	html += `
<tr height="33.58">
<td class="xl83" height="33.58" colspan="16">运维工作遇到的主要问题</td>
</tr>

<tr height="40">
<td class="xl71">编号</td>
<td class="xl72" colspan="11">1.本周运维工作遇到的主要问题</td>
<td class="xl99">是否解决</td>
<td class="xl71">负责人</td>
<td class="xl71">辅助人</td>
<td class="xl71">是否需要协助<br/><font class="font30">(下拉菜单选择)</font></td>
</tr>

<tr height="43.33">
<td class="xl75">1</td>
<td class="xl76" colspan="11">无重大问题</td>
<td class="xl100">已解决</td>
<td class="xl93">熊棋</td>
<td class="xl79">余湘</td>
<td class="xl92">否</td>
</tr>

<tr height="35">
<td class="xl75">2</td>
<td class="xl76" colspan="11"></td>
<td class="xl100"></td>
<td class="xl93"></td>
<td class="xl93"></td>
<td class="xl92"></td>
</tr>

<tr height="35">
<td class="xl75">3</td>
<td class="xl76" colspan="11"></td>
<td class="xl100"></td>
<td class="xl93"></td>
<td class="xl93"></td>
<td class="xl92"></td>
</tr>

<tr height="46.25">
<td class="xl71">编号</td>
<td class="xl72" colspan="12">2.下周可能会遇到的主要问题</td>
<td class="xl72" colspan="2">负责人</td>
<td class="xl71">是否需要协助<br/><font class="font30">(下拉菜单选择)</font></td>
</tr>

<tr height="35">
<td class="xl75">1</td>
<td class="xl76" colspan="12">无预期重大问题</td>
<td class="xl93" colspan="2">熊棋</td>
<td class="xl92">否</td>
</tr>

<tr height="35">
<td class="xl75">2</td>
<td class="xl76" colspan="12"></td>
<td class="xl93" colspan="2"></td>
<td class="xl92"></td>
</tr>

<tr height="35">
<td class="xl75">3</td>
<td class="xl76" colspan="12"></td>
<td class="xl93" colspan="2"></td>
<td class="xl92"></td>
</tr>

<tr height="29">
<td class="xl83" colspan="16">意见及建议</td>
</tr>

<tr height="24">
<td class="xl71">编号</td>
<td class="xl72" colspan="14">意见及建议内容</td>
<td class="xl71">提出人</td>
</tr>

<tr height="24">
<td class="xl75">1</td>
<td class="xl76" colspan="14"></td>
<td class="xl93"></td>
</tr>

<tr height="24">
<td class="xl75">2</td>
<td class="xl76" colspan="14"></td>
<td class="xl93"></td>
</tr>

<tr height="24">
<td class="xl75">3</td>
<td class="xl76" colspan="14"></td>
<td class="xl93"></td>
</tr>
</table>
</body>
</html>`

	return fmt.Sprintf(html, now.Year(), weekNum)
}

func getThisMonday(now time.Time) time.Time {
	weekday := now.Weekday()
	daysSinceMonday := int(weekday) - 1
	if daysSinceMonday < 0 {
		daysSinceMonday = 6 // 周日的情况
	}
	return time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
}

// 是否需要项目组协助
func getAssistanceNeeded(task Task) string {
	if task.Priority >= 2 { // 优先级设置
		return "否"
	}
	return "需要项目组协助"
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
		return "超期未完成" // 必须返回字符串
	}
	return "进行中"
}

func getTaskStatusClass(task Task) string {
	if task.Done {
		return "completed"
	}
	if time.Now().After(task.EndTime) {
		return "overdue" // 必须返回字符串
	}
	return "pending"
}

// 辅助函数：获取当前是第几周（ISO 8601 标准）
func getWeekOfYear() int {
	_, week := time.Now().ISOWeek()
	return week
}
