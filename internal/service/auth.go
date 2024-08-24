package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/SavelyDev/crud-app/pkg/hash"
	"github.com/golang-jwt/jwt"
)

type AuthRepo interface {
	CreateUser(user domain.User) (int, error)
	GetUserId(email, password string) (int, error)
}

const (
	salt       = "rfertertert"
	signingKey = "wqeqweqwe"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	hash := hash.NewSHA1Hasher(salt)

	passwordHash, err := hash.Hash(user.PasswordHash)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = passwordHash
	user.Registered = time.Now()
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	hash := hash.NewSHA1Hasher(salt)

	passwordHash, err := hash.Hash(password)
	if err != nil {
		return "", err
	}

	userId, err := s.repo.GetUserId(email, passwordHash)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(userId),
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix()})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accesToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesToken, &jwt.StandardClaims{}, func(accesToken *jwt.Token) (interface{}, error) {
		if _, ok := accesToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *jwt.StandardClaims")
	}

	id, err := strconv.Atoi(claims.Subject)
	if !ok {
		return 0, err
	}

	return id, nil
}
