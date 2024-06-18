package main

import (
	"LongIM/router"
	"LongIM/utils"
	"fmt"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	//test.TestGorm()
	//fmt.Println("")
	r := router.Router()
	err := r.Run()
	if err != nil {
		fmt.Println(err)
	}
}
