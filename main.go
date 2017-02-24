package main

import (
	_ "linkworld/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {

	logs.EnableFuncCallDepth(true)

	logs.GetLogger("ORM").Println("this is a message of orm")

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
