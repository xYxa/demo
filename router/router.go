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
		user.GET("/tasks", controllers.UserController{}.GetTasks)
		user.POST("/tasks", controllers.UserController{}.CreateTask)
		user.PUT("tasks/:id", controllers.UserController{}.UpdateTask)
		user.DELETE("tasks/:id", controllers.UserController{}.DeleteTask)
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
		order.POST("list", controllers.OrderController{}.GetList)
	}

	return r
}

// Cors 跨域cors部分
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
