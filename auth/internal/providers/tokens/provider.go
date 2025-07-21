package tokens

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Repository interface{}

type TokensProvider struct {
	repository Repository
}

func New(repository Repository) TokensProvider {
	return TokensProvider{
		repository: repository,
	}
}

var (
	accessPrivateKey  string
	accessTokenTTL    = 15 * time.Minute
	refreshPrivateKey string
	refreshTokenTTL   = 7 * 24 * time.Hour
)

func InitJWT(accessKey, refreshKey string) error {
	accessPrivateKey = accessKey
	refreshPrivateKey = refreshKey
	return nil
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Auth tokens...

func (p *TokensProvider) GenerateAccessToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(accessPrivateKey))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (p *TokensProvider) GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(refreshPrivateKey))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (p *TokensProvider) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshPrivateKey), nil
	})

	log.Default().Println("[REFRESH TOKEN]:", token)

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// License keys(tokens)...

func (p *TokensProvider) GenerateRoleKey(
	role string,
	groupID string,
	universityID string,
	enrollmentYear int,
	degree string,
) (string, error) {
	claims := jwt.MapClaims{
		"role":          role,
		"university_id": universityID,
	}

	switch role {
	case "student":
		claims["group_id"] = groupID
		claims["enrollment_year"] = enrollmentYear
	case "teacher":
		claims["degree"] = degree
	default:
		return "", ErrInvalidRole
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshPrivateKey))
}

func (p *TokensProvider) ValidateRoleKey(key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshPrivateKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("key parsing failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidKey
}
