package tokens

import (
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

func (p TokensProvider) GenerateAccessToken(userID string) (string, error) {
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

func (p TokensProvider) GenerateRefreshToken(userID string) (string, error) {
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
