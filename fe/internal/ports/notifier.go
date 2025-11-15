package ports

import "context"

type EmailSender interface {
	SendOTP(ctx context.Context, email, code string) error
}
