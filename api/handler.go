package api

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h Handler) singIn(ctx *gin.Context) {

}

func (h Handler) singUp(ctx *gin.Context) {

}

func (h Handler) create(ctx *gin.Context) {

}

func (h Handler) delete(ctx *gin.Context) {

}

func (h Handler) capsule(ctx *gin.Context) {

}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.singIn)
		auth.POST("/sing-up", h.singUp)
	}

	api := router.Group("/api")
	{
		api.POST("/create", h.create)
		api.DELETE("/delete", h.delete)
		api.GET("/capsule", h.capsule)
	}

	return router
}
