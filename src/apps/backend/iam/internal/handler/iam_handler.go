package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/service"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/iam"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type IamHandler struct {
	svc       service.IamService
	jwtSecret []byte
	pb.UnimplementedIamServiceServer
}

func NewIamHandler(svc service.IamService, jwtSecret string) *IamHandler {
	return &IamHandler{svc: svc, jwtSecret: []byte(jwtSecret)}
}

func (h *IamHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return h.svc.Register(ctx, req)
}

func (h *IamHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return h.svc.Login(ctx, req)
}

func (h *IamHandler) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userID, err := h.getUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	return h.svc.GetUserInfo(ctx, userID)
}

func (h *IamHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	return h.svc.CheckPermission(ctx, req)
}

func (h *IamHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	return h.svc.VerifyToken(ctx, req.Token)
}

func (h *IamHandler) getUserIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}
