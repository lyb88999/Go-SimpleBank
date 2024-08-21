package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)

	username := util.RandString(6)
	role := util.DepositorRole
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, _, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.Equal(t, payload.Role, role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)
	token, _, err := maker.CreateToken(util.RandString(6), util.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgoNone(t *testing.T) {
	payload, err := NewPayload(util.RandString(6), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
