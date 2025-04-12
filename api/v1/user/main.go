package user

import "github.com/gin-gonic/gin"

type UserHandler struct {
	router *gin.RouterGroup
}

func NewUserHandler(router *gin.RouterGroup) *UserHandler {
	userRouter := router.Group("/user")
	return &UserHandler{
		router: userRouter,
	}
}

func (h *UserHandler) SetupRoutes() {
	h.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "user routes",
		})
	})
}
