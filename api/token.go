package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshTokens string `json:"refresh_tokens" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// 在RefreshToken的有效期内刷新AccessToken
func (server *Server) renewAccess(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshTokens)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if session.RefreshToken != req.RefreshTokens {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
