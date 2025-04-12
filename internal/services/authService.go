package services

import (
	"chillfix/internal/repository/user"
	"chillfix/models"
	"context"
	"encoding/json"
	"fmt"

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
}

func NewGoogleAuthService(googleClientID string, googleClientSecret string, googleRedirectURI string, userRepository user.UserRepository) *GoogleAuthService {
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
	}
}

func (s *GoogleAuthService) Login() string {
	url := s.oauth2Config.AuthCodeURL(s.oauthStateString)
	return url
}

func (s *GoogleAuthService) Callback(code string, state string) (string, error) {
	// Check if the state is valid
	if state != s.oauthStateString {
		return "", fmt.Errorf("invalid state")
	}
	// Exchange the authorization code for a token
	token, err := s.oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}

	// Use the token to get user info from Google
	client := s.oauth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	// Parse the user info and return it
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	email := userInfo["email"].(string)

	var name string
	if userInfo["name"] != nil {
		name = userInfo["name"].(string)
	} else {
		name = email
	}
	user := models.NewUser(name, email, "")
	err = s.userRepository.Create(context.Background(), user)
	if err != nil {
		return "", err
	}

	return user.ID.String(), nil
}
func (s *GoogleAuthService) Logout() {
	// Implement logout logic here
	// This could involve clearing session data or tokens
	// For now, we'll just print a message
	fmt.Println("User logged out")
}
