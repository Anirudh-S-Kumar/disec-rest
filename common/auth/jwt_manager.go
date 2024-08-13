package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	privateKey           *ecdsa.PrivateKey
	publicKey            *ecdsa.PublicKey
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJWTManager(privateKey *ecdsa.PrivateKey, access_duration time.Duration, refresh_duration time.Duration) *JWTManager {
	return &JWTManager{
		privateKey,
		&privateKey.PublicKey,
		access_duration,
		refresh_duration}
}

func TempJWTManager() *JWTManager {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return &JWTManager{
		privateKey:           privateKey,
		publicKey:            &privateKey.PublicKey,
		accessTokenDuration:  time.Minute * 15,
		refreshTokenDuration: time.Hour * 24 * 30,
	}
}

func (manager *JWTManager) Generate(username, role string) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("unable to generate uuid: %v", err)
	}

	claims := UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.accessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(manager.privateKey)
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	var claims *UserClaims

	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return manager.publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %+v", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
