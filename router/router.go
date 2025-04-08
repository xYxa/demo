package router

import (
	"demo/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	//跨域 cors
	r.Use(Cors())
	//user路径
	user := r.Group("/user")
	{
		user.GET("/info/:id", controllers.UserController{}.GetUserInfo)
		user.POST("/add", controllers.UserController{}.GetAdd)
		user.GET("/tasks", controllers.UserController{}.GetTasks)          // 获取任务列表
		user.POST("/tasks", controllers.UserController{}.CreateTask)       // 创建新任务
		user.PUT("/tasks/:id", controllers.UserController{}.UpdateTask)    // 更新任务
		user.DELETE("/tasks/:id", controllers.UserController{}.DeleteTask) // 删除任务
		// 周报功能
		user.GET("/weekly-report", controllers.UserController{}.GenerateWeeklyReport) // 支持 format=json/html 参数
		user.GET("/get", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "123get")
		})
		user.POST("/post", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "postMAny")
		})
	}
	//order路径
	order := r.Group("/order")
	{
		order.GET("list", controllers.OrderController{}.GetDailyTasks)
	}

	return r
}

// 跨域 CORS 配置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")

			// 预检请求直接返回
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		c.Next()
	}
}
