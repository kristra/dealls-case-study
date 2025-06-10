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
		payroll := v1.Group("/payrolls")
		{
			payroll.POST("/:year/:month", handlers.UpsertPayroll)
		}
	}

	r.Run()
}
