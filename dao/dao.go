package dao

import (
	"demo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var (
	Db  *gorm.DB
	err error
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

func init() {
	Db, err = gorm.Open(mysql.Open(config.Mysqldb), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	//Db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	Db.AutoMigrate(&Task{})
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	//sqlDB, err := db.DB()

	//// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	//db.SetMaxIdleConns(10)
	//// SetMaxOpenConns 设置打开数据库连接的最大数量。
	//db.SetMaxOpenConns(100)
	//
	//// SetConnMaxLifetime 设置了连接可复用的最大时间。
	//db.SetConnMaxLifetime(time.Hour)
}
