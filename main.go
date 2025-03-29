package main

import (
	"demo/router"
)

func main() {

	r := router.Router()

	r.Run(":8080")
}
