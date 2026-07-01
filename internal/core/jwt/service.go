package corejwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	accessSecret []byte
	accessTTL    time.Duration

	refreshByteLen int
	refreshTTL     time.Duration
}

func NewService(config Config) *Service {
	return &Service{
		accessSecret:   []byte(config.AccessSecret),
		accessTTL:      config.AccessTTL,
		refreshByteLen: config.RefreshByteLen,
		refreshTTL:     config.RefreshTTL,
	}
}

type Pair struct {
	AccessToken  string
	RefreshToken RefreshToken
}

func (s *Service) generateRefresh() (RefreshToken, error) {
	randBytes := make([]byte, s.refreshByteLen)
	if _, err := rand.Read(randBytes); err != nil {
		return RefreshToken{}, err
	}

	str := base64.RawURLEncoding.EncodeToString(randBytes)

	now := time.Now()
	token := RefreshToken{
		Token:     str,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.refreshTTL),
	}

	return token, nil
}

func (s *Service) GenerateAccess(
	role string,
	userId int64,
) (string, error) {
	now := time.Now()

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userId, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := access.SignedString(s.accessSecret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (s *Service) IssuePair(
	role string,
	userId int64,
) (Pair, error) {
	refresh, err := s.generateRefresh()
	if err != nil {
		return Pair{}, fmt.Errorf("issue refresh token: %w", err)
	}

	access, err := s.GenerateAccess(role, userId)
	if err != nil {
		return Pair{}, fmt.Errorf("issue access token: %w", err)
	}

	pair := Pair{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	return pair, nil
}

func (s *Service) ParseAccessToken(access string) (Claims, error) {
	claims := Claims{}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return s.accessSecret, nil
	}

	if _, err := jwt.ParseWithClaims(access, &claims, keyFunc); err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}
