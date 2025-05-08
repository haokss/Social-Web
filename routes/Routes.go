package routes

import (
	"todo_list/api"
	"todo_list/middleware"
	sse "todo_list/package/SSE"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 明确指定允许的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 预检请求缓存 24 小时
	}))

	// 初始化sse通知机制
	broker := sse.NewBroker()
	defer broker.Shutdown()

	// 处理 OPTIONS 请求
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Set("sseBroker", broker)
		c.Next()
	})

	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))

	v1 := r.Group("api/v1")
	{
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		// TODO: 1、忘记密码 2、修改密码
		// 新建认证分组
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			// 用户信息
			authed.GET("user_info", api.UserInfo)
			authed.PUT("user_info", api.UpdateUserInfo)
			// 备忘录任务
			authed.POST("task", api.CreateTask)
			authed.GET("task/:id", api.ShowTask)
			authed.POST("tasks", api.ShowTaskAll)
			authed.PUT("task/:id", api.UpdateTask)
			authed.POST("task/search", api.SearchTask)
			authed.DELETE("task/:id", api.DeleteTask)
			// 定时提醒任务
			authed.POST("timing_task", api.CreateTimingTask)
			// authed.GET("timing_task/:id", api.ShowTimingTask)
			authed.POST("timing_tasks", api.ShowTimingTaskAll)
			authed.PUT("timing_task/:id", api.UpdateTimingTask)
			// authed.POST("timing_task/search", api.SearchTimingTask)
			authed.POST("timing_task/delete_batch", api.DeleteTimingTasks)

			// 任务活动统计
			authed.GET("tasks/stats", api.GetTaskStats)
			authed.GET("tasks/trend", api.GetTaskTrend)
			authed.GET("tasks/type-distribution", api.GetTaskTypeDistribution)
			authed.GET("tasks/upcoming", api.ShowUpcomingTasks)
			authed.GET("/tasks/high-priority", api.ShowHighPriorityTasks)

			// 亲属关系
			authed.POST("/relative_info", api.CreateRelative)
			authed.PUT("/relative_info/:id", api.UpdateRelative)
			authed.GET("/relative_info", api.GetAllRelatives)
			authed.DELETE("/relative_info", api.DeleteRelative)

			// 同事关系
			authed.POST("/colleague", api.CreateColleague)
			authed.PUT("/colleague/:id", api.UpdateColleague)
			authed.DELETE("/colleague/batch", api.BatchDeleteColleagues)
			authed.GET("/colleagues", api.ListTimingTask)

			// 朋友管理
			authed.POST("/friend", api.CreateFriend)
			authed.GET("/friend/list", api.ListFriend)
			authed.PUT("/friend/:id", api.UpdateFriend)
			authed.DELETE("/friend/batch_delete", api.BatchDeleteFriends)
			// 同学管理
			authed.POST("/classmate", api.CreateClassmate)
			authed.GET("/classmate/list", api.ListClassmate)
			authed.PUT("/classmate/:id", api.UpdateClassmate)
			// authed.DELETE("/friends/:id", api.BatchDeleteClassmate) // 单个删除你已有
			authed.DELETE("/classmate/batch_delete", api.BatchDeleteClassmate)

			// 关系地图
			authed.GET("/relative/unset_map_relative", api.GetUnsetMapRelatives)
			authed.GET("/relative/unset_map_colleague", api.GetUnsetMapColleagues)
			authed.GET("/relative/unset_map_friend", api.GetUnsetMapFriends)
			authed.GET("/relative/unset_map_classmate", api.GetUnsetMapClassmates)

			authed.POST("/relative/point", api.CreatePoint)
			authed.GET("/relative/point", api.ListPoints)
			authed.PUT("/relative/point", api.UpdatePoint)
			authed.DELETE("/relative/point/:id", api.DeletePoint)

			// 账单导入GET
			authed.POST("/import-bill", api.ImportBill)
			authed.GET("/bills", api.GetBillList)

			// sse消息推送
			authed.GET("/push/notifications", middleware.SSEAuth(), sse.SSEHandler())

			// 图片上传
			authed.POST("/upload/ava", api.UploadAva)
			authed.POST("/upload/image", api.UploadImage)
			authed.GET("/upload/image", api.GetImage)
		}
	}
	return r
}
