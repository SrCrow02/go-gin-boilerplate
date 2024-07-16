package router

import (
	"github.com/gin-gonic/gin"
	"template/handler"
)

func Routes(port string) {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/add", handler.Add)
		v1.GET("/all", handler.GetAll)
		v1.PATCH("/update/:id", handler.Update)
		v1.DELETE("/delete/:id", handler.Delete)
	}

	router.Run("localhost:" + port)
}
