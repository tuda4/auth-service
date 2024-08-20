package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/tuda4/mb-backend/db/sqlc"
	"github.com/tuda4/mb-backend/util"
)

type PayloadVerifyEmail struct {
	Email string `json:"email"`
}

const (
	TaskVerifyEmail = "task:send_verify_email"
)

func (d *RedisTaskDistributor) DistributeTaskVerifyEmail(
	ctx context.Context,
	payload *PayloadVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	fmt.Print(jsonPayload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %v", err)
	}
	task := asynq.NewTask(TaskVerifyEmail, jsonPayload, opts...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}

	log.Info().Str("type", info.Type).Bytes("payload", info.Payload).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueue task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessorSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %v", err)
	}

	account, err := processor.store.GetProfileAccount(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("could not get user: %v", err)
	}
	verifyCode := util.RandomString(32)

	_, err = processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		AccountID:  account.AccountID,
		SecretCode: verifyCode,
		IsUsed:     false,
	})
	if err != nil {
		return fmt.Errorf("could not create verify email: %v", err)
	}

	subject := `Tuda4 test send email verify account`
	verifyURL := fmt.Sprintf("http://localhost:3000/?account_id=%v&secret_code=%v", account.AccountID, verifyCode)
	content := fmt.Sprintf("Dear, %s"+verifyURL, account.FirstName.String)
	to := []string{account.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", account.Email).Msg("processed task")

	return nil
}
