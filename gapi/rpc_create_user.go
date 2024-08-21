package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"github.com/lyb88999/Go-SimpleBank/pb"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/lyb88999/Go-SimpleBank/val"
	"github.com/lyb88999/Go-SimpleBank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		err := invalidArgumentError(violations)
		return nil, err
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(time.Second * 10),
				asynq.Queue(worker.QueueCritical),
			}
			err = server.taskDistributor.DistributeSendVerifyEmail(ctx, taskPayload, opts...)
			return err
		},
	}
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			if db.ErrorConstraintName(err) == "users_pkey" {
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err.Error())
			}
			if db.ErrorConstraintName(err) == "users_email_key" {
				return nil, status.Errorf(codes.AlreadyExists, "email address already exists: %s", err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// TODO: use db transaction

	rsp := &pb.CreateUserResponse{
		User: ConvertUser(txResult.User),
	}
	return rsp, err
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateUsername(req.GetUsername())
	if err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	err = val.ValidatePassword(req.GetPassword())
	if err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	err = val.ValidateEmail(req.GetEmail())
	if err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	err = val.ValidateFullName(req.GetFullName())
	if err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	return
}
