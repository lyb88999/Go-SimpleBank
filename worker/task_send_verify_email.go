package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	db "github.com/lyb88999/Go-SimpleBank/db/sqlc"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("cannot enquene the task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", jsonPayload).
		Str("queue", info.Queue).Int("max_entry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload *PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", asynq.SkipRetry)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		//if errors.Is(err, ErrRecordNotFound) {
		//	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		//}
		// 注释以后失败后会重试
		return fmt.Errorf("cannot get user:%w", err)
	}

	// TODO: send email to user
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}
	subject := "Welcome to Simple Bank"
	verifyURL := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you registering with us!<br/>
	Please <a href="%s">Clike here</a> to verify your email address.</br>`, user.FullName, verifyURL)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email:%w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}
