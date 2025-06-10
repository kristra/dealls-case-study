package route

import (
	"dealls-case-study/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}
		payroll := v1.Group("/payrolls")
		{
			payroll.POST("/:year/:month", handlers.UpsertPayroll)
		}
	}

	r.Run()
}
