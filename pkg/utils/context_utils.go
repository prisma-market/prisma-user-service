package utils

import (
	"context"
	"fmt"
)

type contextKey string

const (
	userClaimsKey contextKey = "user_claims"
)

// SetUserContext 컨텍스트에 JWT claims 저장
func SetUserContext(ctx context.Context, claims *JWTClaim) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserFromContext 컨텍스트에서 JWT claims 조회
func GetUserFromContext(ctx context.Context) (*JWTClaim, error) {
	claims, ok := ctx.Value(userClaimsKey).(*JWTClaim)
	if !ok {
		return nil, fmt.Errorf("user claims not found in context")
	}
	return claims, nil
}
