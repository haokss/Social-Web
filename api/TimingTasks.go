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
		res := create_timing_task.Create(c, claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 展示用户定时活动
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

// 展示用户所有定时活动
func ShowTimingTaskAll(c *gin.Context) {
	var timing_tasks service.ShowTimingTaskAllService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&timing_tasks); err == nil {
		res := timing_tasks.ShowAll(c, claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

// 更新一个备忘录
func UpdateTimingTask(c *gin.Context) {
	var task service.UpdateTimingTaskService
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
func DeleteTimingTasks(c *gin.Context) {
	var service service.DeleteTimingTasksService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Delete(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
