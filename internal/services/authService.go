package services

import (
	"chillfix/internal/repository/user"
	"chillfix/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthService struct {
	googleClientID     string
	googleClientSecret string
	googleRedirectURI  string
	oauth2Config       oauth2.Config
	oauthStateString   string
	userRepository     user.UserRepository
	tokenService       *TokenService
}

func NewGoogleAuthService(googleClientID string, googleClientSecret string, googleRedirectURI string, userRepository user.UserRepository, tokenService *TokenService) *GoogleAuthService {
	oauth2State := "randomstate"
	oauthStateString := oauth2State
	googleScopes := []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	oauth2Config := oauth2.Config{
		ClientID:     googleClientID,     // Use the provided Google client ID
		ClientSecret: googleClientSecret, // Use the provided Google client secret
		RedirectURL:  googleRedirectURI,  // Use the provided redirect URI
		Scopes:       googleScopes,
		Endpoint:     google.Endpoint,
	}
	return &GoogleAuthService{
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		googleRedirectURI:  googleRedirectURI,
		userRepository:     userRepository,
		oauth2Config:       oauth2Config,
		oauthStateString:   oauthStateString,
		tokenService:       tokenService,
	}
}

func (s *GoogleAuthService) Login() string {
	url := s.oauth2Config.AuthCodeURL(s.oauthStateString)
	return url
}

func (s *GoogleAuthService) Callback(code string, state string) (accessToken, refreshToken string, userId string, err error) {
	// Check if the state is valid
	if state != s.oauthStateString {
		return "", "", "", fmt.Errorf("invalid state")
	}
	// Exchange the authorization code for a token
	token, err := s.oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("Error exchanging code for token:", err)
		return "", "", "", err
	}

	// Use the token to get user info from Google
	client := s.oauth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Println("Error getting user info:", err)
		return "", "", "", err
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != 200 {
		log.Println("Error getting user info:", resp.Status)
		return "", "", "", fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	// Parse the user info and return it
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Println("Error decoding user info:", err)
		return "", "", "", err
	}

	email := userInfo["email"].(string)

	var name string
	if userInfo["name"] != nil {
		name = userInfo["name"].(string)
	} else {
		name = email
	}

	var user *models.User
	// Check if the user already exists in the database

	user, err = s.userRepository.FindByEmail(context.Background(), email)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			user = models.NewUser(name, email, "")
			err = s.userRepository.Create(context.Background(), user)
			if err != nil {
				log.Println("Error creating user:", err)
				return "", "", "", err
			}
		} else {
			log.Println("Error finding user:", err)
			return "", "", "", err
		}
	}

	// Generate a access and refresh token
	accessToken, err = s.tokenService.GenerateToken(map[string]interface{}{
		"email":  email,
		"name":   name,
		"userid": user.Id,
	}, 15)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err = s.tokenService.GenerateToken(map[string]interface{}{
		"email":  email,
		"name":   name,
		"userid": user.Id,
	}, 60*24*7) // 7 days
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", "", err
	}

	// Save the tokens to the user
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	err = s.userRepository.Update(context.Background(), user)
	if err != nil {
		log.Println("Error updating user:", err)
		return "", "", "", err
	}
	// return access and refresh token and userid
	return accessToken, refreshToken, user.Id, nil
}
func (s *GoogleAuthService) Logout() {
	// Implement logout logic here
	// This could involve clearing session data or tokens
	// For now, we'll just print a message
	fmt.Println("User logged out")
}
