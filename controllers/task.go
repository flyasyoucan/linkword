package controllers

import (
	"encoding/json"
	"linkworld/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type TaskController struct {
	beego.Controller
}

type taskinfo struct {
	Tid string
}

type taskRes struct {
	RespCode  int
	Resultmsg string
	Data      taskinfo
}

// @Title AddTask
// @Description add task
// @Param	body		body 	models.User	true		"body for task content"
// @Success 200 {int} models.Task.Id
// @Failure 403 body is empty
// @router / [post]
func (t *TaskController) Post() {

	var task models.Task
	json.Unmarshal(t.Ctx.Input.RequestBody, &task)
	tid := models.AddTask(task)

	var res taskRes

	res.RespCode = 0
	res.Resultmsg = "success"
	res.Data.Tid = tid

	t.Data["json"] = &res

	t.ServeJSON()
}

// @Title Add
// @Description add a task
// @Param	body		body 	models.Task	true		"body for user content"
// @Success 200 {string} models.Task.Id
// @Failure 403 body is empty
// @router /add [post]
func (t *TaskController) Add() {

	var task models.Task
	json.Unmarshal(t.Ctx.Input.RequestBody, &task)
	tid := models.AddTask(task)

	logs.Debug("add new task:")
	var res taskRes

	res.RespCode = 0
	res.Resultmsg = "success"
	res.Data.Tid = tid

	t.Data["json"] = &res

	t.ServeJSON()
}
