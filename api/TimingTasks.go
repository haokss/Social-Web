package api

import (
	"todo_list/package/utils"
	"todo_list/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

// 创建一个定时提醒任务
func CreateTimingTask(c *gin.Context) {
	var create_timing_task service.CreateTimingTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&create_timing_task); err == nil {
		res := create_timing_task.Create(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 展示用户的备忘录
func ShowTimingTask(c *gin.Context) {
	var task service.ShowTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&task); err == nil {
		res := task.Show(claim.Id, c.Param("id"))
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 展示用户所有的备忘录
func ShowTimingTaskAll(c *gin.Context) {
	var timing_tasks service.ShowTimingTaskAllService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&timing_tasks); err == nil {
		res := timing_tasks.ShowAll(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 更新一个备忘录
func UpdateTimingTask(c *gin.Context) {
	var task service.UpdateTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&task); err == nil {
		res := task.Update(claim.Id, c.Param("id"))
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 查询备忘录
func SearchTimingTask(c *gin.Context) {
	var tasks service.SearchTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&tasks); err == nil {
		res := tasks.Search(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 删除备忘录
func DeleteTimingTask(c *gin.Context) {
	var delete_task service.DeleteTaskService
	// claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&delete_task); err == nil {
		res := delete_task.Delete(c.Param("id"))
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
