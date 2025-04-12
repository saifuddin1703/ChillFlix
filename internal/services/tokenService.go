package services

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	// Add any dependencies you need here, such as a database client or a JWT library
	jwtSecret string
}

func NewTokenService(jwtSecret string) *TokenService {
	return &TokenService{
		jwtSecret: jwtSecret,
	}
}

// GenerateToken generates a JWT token for the given user ID
func (s *TokenService) GenerateToken(payload map[string]interface{}, expiry int) (string, error) {
	// Implement JWT token generation logic here
	payload["exp"] = time.Now().Add(time.Duration(expiry) * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateToken validates the given JWT token and returns the user ID if valid
func (s *TokenService) ValidateToken(tokenString string) (string, error) {
	// Implement JWT token validation logic here
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userid"].(string)
		return userID, nil
	}
	return "", jwt.ErrInvalidKey
}
