package grpcapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/fajaramaulana/simple_bank_project/db/sqlc"
	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/token"
	"github.com/fajaramaulana/simple_bank_project/pb"
	"github.com/fajaramaulana/simple_bank_project/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	detailLogin, err := s.db.GetDetailLoginByUsername(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			// return nil, ErrUserNotFound
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	// check password
	err = util.CheckPasswordBcrypt(req.GetPassword(), detailLogin.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	// generate token with paseto
	maker, err := token.NewPasetoMaker(s.config.TokenSymmetricKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create token maker: %v", err)
	}
	accessTokenDuration, err := time.ParseDuration(s.config.AccessTokenDuration.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse access token duration: %v", err)
	}

	role := detailLogin.Role

	accessToken, _, err := maker.CreateToken(detailLogin.UserUuid.String(), accessTokenDuration, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshTokenDuration, err := time.ParseDuration(s.config.RefreshTokenDuration.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse refresh token duration: %v", err)
	}

	refreshToken, payloadRefresh, err := maker.CreateToken(detailLogin.UserUuid.String(), refreshTokenDuration, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	arg := db.CreateSessionParams{
		ID:           payloadRefresh.ID,
		UserUuid:     detailLogin.UserUuid,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    payloadRefresh.ExpiredAt,
	}

	// save refresh token to db
	session, err := s.db.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	res := &pb.LoginUserResponse{
		SessionId:    session.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			UserUuid:  detailLogin.UserUuid.String(),
			Username:  detailLogin.Username,
			FullName:  detailLogin.FullName,
			Email:     detailLogin.Email,
			CreatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		},
	}

	return res, nil
}
