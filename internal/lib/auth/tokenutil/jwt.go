package tokenutil

import (
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type tokenRepository interface {
	BlacklistTokenByID(tokenStr string) error
	IsBlacklisted(tokenStr string) (bool, error)
}

type JWTTokenManager struct {
	AccessSecret    string
	RefreshSecret   string
	AccessExpiry    time.Duration
	RefreshExpiry   time.Duration
	TokenRepository tokenRepository
}

func NewJWTTokenManager(
	accessSecret string,
	refreshSecret string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
	tokenRepository tokenRepository,
) *JWTTokenManager {
	return &JWTTokenManager{
		AccessSecret:    accessSecret,
		RefreshSecret:   refreshSecret,
		AccessExpiry:    accessExpiry,
		RefreshExpiry:   refreshExpiry,
		TokenRepository: tokenRepository,
	}
}

// CreateUserAccessToken creates a new access token for the user.
func (tm *JWTTokenManager) CreateUserAccessToken(user *models.User) (string, error) {
	return tm.createAccessToken(user.ID)
}

// CreateUserRefreshToken creates a new refresh token for the user.
func (tm *JWTTokenManager) CreateUserRefreshToken(user *models.User) (string, error) {
	return tm.createRefreshToken(user.ID)
}

func (tm *JWTTokenManager) createAccessToken(id int64) (string, error) {
	return tm.createToken(id, tm.AccessSecret, tm.AccessExpiry)
}

func (tm *JWTTokenManager) createRefreshToken(id int64) (string, error) {
	return tm.createToken(id, tm.RefreshSecret, tm.RefreshExpiry)
}

func (tm *JWTTokenManager) createToken(id int64, secret string, expiry time.Duration) (string, error) {
	const op = "tokenutil.createToken"

	now := time.Now()
	exp := now.Add(expiry)
	claims := &jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   strconv.FormatInt(id, 10),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedToken, nil
}
