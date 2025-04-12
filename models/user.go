package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"-"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUser(username string, email string, password string) *User {
	user := &User{
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if user.Password != "" {
		hashedPassword, err := user.HashPassword()
		if err != nil {
			return nil
		}
		user.Password = hashedPassword
	}
	return user
}

func (u *User) SetID(id primitive.ObjectID) {
	u.ID = id
}

func (u *User) HashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	u.Password = string(hashedPassword)
	return string(hashedPassword), nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) ValidatePassword(password string) bool {
	return u.Password == password
}
