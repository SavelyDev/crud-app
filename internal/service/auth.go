package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/golang-jwt/jwt"
)

type AuthRepo interface {
	CreateUser(user domain.User) (int, error)
	GetUserId(email, password string) (int, error)
}

type TokensRepo interface {
	CreateSession(refreshToken domain.RefreshSession) error
	GetSession(refreshToken string) (domain.RefreshSession, error)
}

type PasswordHash interface {
	Hash(password string) (string, error)
}

type AuthService struct {
	repo       AuthRepo
	tokensRepo TokensRepo
	hasher     PasswordHash
	tokenTTL   time.Duration
	signingKey []byte
}

func NewAuthService(repo AuthRepo, tokensRepo TokensRepo, hasher PasswordHash,
	tokenTTL time.Duration, signingKey []byte) *AuthService {
	return &AuthService{repo: repo, tokensRepo: tokensRepo, hasher: hasher,
		tokenTTL: tokenTTL, signingKey: signingKey}
}

func (s *AuthService) SignUp(user domain.User) (int, error) {
	passwordHash, err := s.hasher.Hash(user.PasswordHash)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = passwordHash
	user.Registered = time.Now()
	return s.repo.CreateUser(user)
}

func (s *AuthService) SignIn(user domain.SignInInput) (string, string, error) {
	passwordHash, err := s.hasher.Hash(user.Password)
	if err != nil {
		return "", "", err
	}

	userId, err := s.repo.GetUserId(user.Email, passwordHash)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(userId)
}

func (s *AuthService) generateTokens(userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(userId),
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix()})

	accesToken, err := t.SignedString(s.signingKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.tokensRepo.CreateSession(domain.RefreshSession{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30),
	}); err != nil {
		return "", "", err
	}

	return accesToken, refreshToken, nil
}

func (s *AuthService) ParseToken(accesToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesToken, &jwt.StandardClaims{}, func(accesToken *jwt.Token) (interface{}, error) {
		if _, ok := accesToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.signingKey, nil
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

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := s.tokensRepo.GetSession(refreshToken)
	if err != nil {
		return "", "", err
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")

	}

	return s.generateTokens(token.UserId)
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)
	c := rand.NewSource(time.Now().Unix())
	r := rand.New(c)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
