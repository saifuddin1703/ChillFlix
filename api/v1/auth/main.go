package auth

import (
	"chillfix/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	googleAuthService *services.GoogleAuthService
	router            *gin.RouterGroup
}

func NewAuthHandler(router *gin.RouterGroup, googleAuthService *services.GoogleAuthService) *AuthHandler {
	authRouter := router.Group("/auth")
	return &AuthHandler{
		googleAuthService: googleAuthService,
		router:            authRouter,
	}
}

func (h *AuthHandler) SetupRoutes() {
	// Login route
	h.router.POST("/login", h.login)

	// Register route
	h.router.POST("/register", h.register)

	// Google login route
	h.router.GET("/google/login", func(ctx *gin.Context) {
		url := h.googleAuthService.Login()

		ctx.Redirect(http.StatusTemporaryRedirect, url)
	})

	// OAuth2.0 callback route
	h.router.GET("/google/callback", h.oauthCallback)
}

// create oauth2.0 callback
func (h *AuthHandler) oauthCallback(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	// Exchange the authorization code for a token
	userId, err := h.googleAuthService.Callback(code, state)
	fmt.Println("userId", userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	// ✅ Generate your own access/refresh tokens
	// accessToken, err := generateJWT(userInfo.Email, time.Minute*15)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
	// 	return
	// }

	// refreshToken, err := generateJWT(userInfo.Email, time.Hour*24*7)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
		"user":          userId,
	})
}
func (h *AuthHandler) login(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "login endpoint",
	})
}

func (h *AuthHandler) register(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "register endpoint",
	})
}
