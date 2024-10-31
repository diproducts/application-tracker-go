package tokenutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// TODO: decide what to do with token repository
type tokenRepository interface {
	BlacklistTokenByID(ctx context.Context, userID int64, tokenID string, expiry time.Time) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTTokenManager struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewJWTTokenManager(
	accessSecret string,
	refreshSecret string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
) *JWTTokenManager {
	return &JWTTokenManager{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
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

// ExtractUserIDFromAccessToken extracts user id from access token.
func (tm *JWTTokenManager) ExtractUserIDFromAccessToken(tokenStr string) (int64, error) {
	return tm.extractUserIDFromToken(tokenStr, tm.AccessSecret)
}

// ExtractUserIDFromRefreshToken extracts user id from refresh token.
func (tm *JWTTokenManager) ExtractUserIDFromRefreshToken(tokenStr string) (int64, error) {
	return tm.extractUserIDFromToken(tokenStr, tm.RefreshSecret)
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

func (tm *JWTTokenManager) parseToken(tokenStr, secret string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return token, err
}

func (tm *JWTTokenManager) extractUserIDFromToken(tokenStr, secret string) (int64, error) {
	idStr, err := tm.extractSubjectFromToken(tokenStr, secret)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}

	return id, nil
}

func (tm *JWTTokenManager) extractSubjectFromToken(tokenStr, secret string) (string, error) {
	const op = "tokenutil.extractSubjectFromToken"

	claims := &jwt.RegisteredClaims{}
	token, err := tm.parseToken(tokenStr, secret, claims)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}
