package gapi

import (
	"context"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"github.com/lyb88999/Go-SimpleBank/pb"
	"github.com/lyb88999/Go-SimpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		err := invalidArgumentError(violations)
		return nil, err
	}
	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to verify email")
	}
	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, err
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateEmailID(req.GetEmailId())
	if err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	err = val.ValidateSecretCode(req.GetSecretCode())
	if err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}
	return
}
