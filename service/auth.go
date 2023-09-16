package service

import (
	"errors"
	"payint"
	"payint/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("My_key")

type Claims struct {
	*jwt.StandardClaims
	UserName string `json:"username"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateAdmin(admin payint.User) error {
	return s.repo.CreateAdmin(admin)
}

func (s *AuthService) CreateUser(user payint.User) error {
	return s.repo.CreateUser(user)
}

func (s *AuthService) BlockUser(user payint.User) error {
	return s.repo.BlockUser(user)
}

func (s *AuthService) UnBlockUser(user payint.User) error {
	return s.repo.UnBlockUser(user)
}

func (s *AuthService) GetUser(username, password string) (payint.User, error) {
	return s.repo.GetUser(username, password)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	hashedPassword, _ := repository.HashePassword(password)
	user, err := s.repo.GetUser(username, hashedPassword)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserName: user.Name,
	})

	return token.SignedString(JwtKey)
}

func (s *AuthService) ParseToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("Invasid signing method")
		}
		return []byte(JwtKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserName, nil
}
