package gapi

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"github.com/lyb88999/Go-SimpleBank/pb"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/lyb88999/Go-SimpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	accessibleRoles := []string{util.DepositorRole, util.BankerRole}
	authPayload, err := server.authorizeUser(ctx, accessibleRoles)
	if err != nil {
		return nil, unauthorizedError(err)
	}
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		err := invalidArgumentError(violations)
		return nil, err
	}
	if authPayload.Role != util.BankerRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other users' info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}
	rsp := &pb.UpdateUserResponse{
		User: ConvertUser(user),
	}
	return rsp, err
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateUsername(req.GetUsername())
	if err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.GetPassword() != "" {
		err = val.ValidatePassword(req.GetPassword())
		if err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.GetEmail() != "" {
		err = val.ValidateEmail(req.GetEmail())
		if err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.GetFullName() != "" {
		err = val.ValidateFullName(req.GetFullName())
		if err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	return
}
