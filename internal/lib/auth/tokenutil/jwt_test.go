package tokenutil_test

import (
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/diproducts/application-tracker-go/internal/lib/auth/tokenutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	userID        int64 = 7
	accessSecret        = "test_access_secret"
	refreshSecret       = "test_refresh_secret"
	accessTTL           = time.Duration(10 * time.Minute)
	refreshTTL          = time.Duration(7 * 24 * time.Hour) // 1 week
)

func TestJWTTokenManager_CreateUserAccessToken(t *testing.T) {
	tm := tokenutil.NewJWTTokenManager(accessSecret, refreshSecret, accessTTL, refreshTTL)
	user := models.User{ID: userID}
	now := time.Now()
	const deltaSeconds = 1

	tokenStr, err := tm.CreateUserAccessToken(&user)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	claims := &jwt.RegisteredClaims{}
	token, err := parseToken(tokenStr, tm.AccessSecret, claims)
	require.NoError(t, err)
	require.True(t, token.Valid)

	gotUserID, err := tm.ExtractUserIDFromAccessToken(tokenStr)
	assert.NoError(t, err)
	_, err = tm.ExtractUserIDFromAccessToken("err" + tokenStr)
	assert.ErrorIs(t, err, tokenutil.ErrInvalidToken)

	assert.NotEmpty(t, claims.ID)
	assert.InDelta(t, now.Unix(), claims.IssuedAt.Unix(), deltaSeconds)
	assert.InDelta(t, now.Unix(), claims.NotBefore.Unix(), deltaSeconds)
	assert.InDelta(t, now.Add(tm.AccessExpiry).Unix(), claims.ExpiresAt.Unix(), deltaSeconds)
	assert.Equal(t, userID, gotUserID)
}

func TestJWTTokenManager_CreateUserRefreshToken(t *testing.T) {
	tm := tokenutil.NewJWTTokenManager(accessSecret, refreshSecret, accessTTL, refreshTTL)
	user := models.User{ID: userID}
	now := time.Now()
	const deltaSeconds = 1

	tokenStr, err := tm.CreateUserRefreshToken(&user)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	claims := &jwt.RegisteredClaims{}
	token, err := parseToken(tokenStr, tm.RefreshSecret, claims)
	require.NoError(t, err)
	require.True(t, token.Valid)

	gotUserID, err := tm.ExtractUserIDFromRefreshToken(tokenStr)
	assert.NoError(t, err)
	_, err = tm.ExtractUserIDFromRefreshToken("err" + tokenStr)
	assert.ErrorIs(t, err, tokenutil.ErrInvalidToken)

	assert.NotEmpty(t, claims.ID)
	assert.InDelta(t, now.Unix(), claims.IssuedAt.Unix(), deltaSeconds)
	assert.InDelta(t, now.Unix(), claims.NotBefore.Unix(), deltaSeconds)
	assert.InDelta(t, now.Add(tm.RefreshExpiry).Unix(), claims.ExpiresAt.Unix(), deltaSeconds)
	assert.Equal(t, userID, gotUserID)
}

func parseToken(tokenStr, secret string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return token, err
}
