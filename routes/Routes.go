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
			// authed.PUT("timing_task/:id", api.UpdateTimingTask)
			// authed.POST("timing_task/search", api.SearchTimingTask)
			// authed.DELETE("timing_task/:id", api.DeleteTimingTask)

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
