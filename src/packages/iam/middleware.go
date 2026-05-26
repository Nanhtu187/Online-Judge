package iam

import (
	"context"
	"errors"
	"strings"

	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/iam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

type AuthMiddleware struct {
	iamClient   pb.IamServiceClient
	internalKey string
}

func NewAuthMiddleware(iamClient pb.IamServiceClient, internalKey string) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient:   iamClient,
		internalKey: internalKey,
	}
}

func (m *AuthMiddleware) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip auth for specific methods
	if strings.Contains(info.FullMethod, "/Register") || strings.Contains(info.FullMethod, "/Login") {
		return handler(ctx, req)
	}

	// 1. Check for Internal Key first (Service-to-Service)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		internalKeys := md.Get("x-internal-key")
		if len(internalKeys) > 0 && m.internalKey != "" && internalKeys[0] == m.internalKey {
			ctx = context.WithValue(ctx, UserIDKey, "SYSTEM")
			return handler(ctx, req)
		}
	}

	// 2. Fallback to JWT
	userID, err := m.getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ctx = context.WithValue(ctx, UserIDKey, userID)
	return handler(ctx, req)
}

func (m *AuthMiddleware) AuthorizeInterceptor(methodPerms map[string]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		permission, ok := methodPerms[info.FullMethod]
		if !ok {
			// No permission required for this method
			return handler(ctx, req)
		}

		userID := GetUserID(ctx)
		if userID == "" {
			return nil, status.Error(codes.Unauthenticated, "missing user identity")
		}

		// SYSTEM user (Internal API Key) is allowed to perform any internal action
		if userID == "SYSTEM" {
			return handler(ctx, req)
		}

		resp, err := m.iamClient.CheckPermission(ctx, &pb.CheckPermissionRequest{
			UserId:     userID,
			Permission: permission,
			Scope:      pb.Scope_SCOPE_GLOBAL,
		})

		if err != nil {
			return nil, status.Error(codes.Internal, "permission check failed")
		}

		if !resp.Allowed {
			return nil, status.Error(codes.PermissionDenied, "unauthorized")
		}

		return handler(ctx, req)
	}
}

func (m *AuthMiddleware) getUserIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
	
	// Call IAM to verify token
	resp, err := m.iamClient.VerifyToken(ctx, &pb.VerifyTokenRequest{
		Token: tokenString,
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	return resp.UserId, nil
}

func (m *AuthMiddleware) CheckPermission(ctx context.Context, permission string, scope pb.Scope, scopeID string) (bool, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return false, errors.New("unauthenticated")
	}

	resp, err := m.iamClient.CheckPermission(ctx, &pb.CheckPermissionRequest{
		UserId:     userID,
		Permission: permission,
		Scope:      scope,
		ScopeId:    scopeID,
	})
	if err != nil {
		return false, err
	}

	return resp.Allowed, nil
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
