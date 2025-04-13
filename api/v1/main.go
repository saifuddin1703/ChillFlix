package v1

import (
	"chillfix/api/v1/auth"
	"chillfix/api/v1/user"
	"chillfix/api/v1/video"
	"chillfix/config"
	"chillfix/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type V1Handler struct {
	router *gin.RouterGroup
}

func NewV1Handler(router *gin.RouterGroup) *V1Handler {
	v1Router := router.Group("/v1")
	return &V1Handler{
		router: v1Router,
	}
}

func (h *V1Handler) SetupRoutes() {

	h.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "v1 routes",
		})
	})

	// Setup the user routes
	userRepo, err := services.NewUserService()
	if err != nil {
		fmt.Println("Error creating user service:", err)
		return
	}

	// Setup the auth routes
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	tokenService := services.NewTokenService(config.GetJWTSecret())
	authRouter := auth.NewAuthHandler(h.router, services.NewGoogleAuthService(config.GetGoogleClientID(), config.GetGoogleClientSecret(), config.GetGoogleRedirectURI(), userRepo.UserRepository, tokenService))
	authRouter.SetupRoutes()

	// authMiddleware := middleware.NewAuthMiddleware(tokenService)
	// h.router.Use(authMiddleware.Authenticate())
	userRouter := user.NewUserHandler(h.router)
	userRouter.SetupRoutes()
	// Setup the video routes
	videoRouter := video.NewVideoHandler(h.router)
	videoRouter.SetupRoutes()
}
