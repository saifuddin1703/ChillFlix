package api

import (
	v1 "chillfix/api/v1"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type API struct {
	Port   int
	engine *gin.Engine
}

func NewAPI(port int) error {
	api := &API{
		Port:   port,
		engine: gin.Default(),
	}
	api.SetupRoutes()
	err := api.Start()
	if err != nil {
		return err
	}
	return nil
}

func (a *API) Start() error {
	// Start the API server
	err := a.engine.Run(fmt.Sprintf(":%d", a.Port))
	if err != nil {
		return err
	}
	return nil
}

func (a *API) SetupRoutes() {
	a.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Or specify "http://localhost:5500"
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router := a.engine.Group("/api")
	// Setup the routes
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Setup the v1 routes
	v1 := v1.NewV1Handler(router)
	v1.SetupRoutes()
}
