package gapi

import (
	"context"
	"fmt"
	"github.com/lyb88999/Go-SimpleBank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	// split the bearer token
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}
	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token:%s", err)
	}
	if !hashPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}
	return payload, nil
}

func hashPermission(role string, accessibleRoles []string) bool {
	for _, accessibleRole := range accessibleRoles {
		if role == accessibleRole {
			return true
		}
	}
	return false
}
