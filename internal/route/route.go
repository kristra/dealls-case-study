package route

import (
	"dealls-case-study/internal/handlers"
	"dealls-case-study/internal/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/login", handlers.Login)

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.AuthMiddleware())
	{
		attendance := v1.Group("/attendances")
		{
			attendance.POST("/check-in", handlers.CheckInAttendance)
			attendance.POST("/check-out", handlers.CheckOutAttendance)
		}
		payroll := v1.Group("/payrolls")
		{
			payroll.POST("/:year/:month", middlewares.AdminOnly(), handlers.UpsertPayroll)
		}
	}

	r.Run()
}
