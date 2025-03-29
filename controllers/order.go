package controllers

import "github.com/gin-gonic/gin"

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
