package email

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/escalopa/inno-vkode/internal/ports"
)

type LogSender struct {
	log zerolog.Logger
}

func NewLogSender(log zerolog.Logger) *LogSender {
	return &LogSender{log: log}
}

func (s *LogSender) SendOTP(ctx context.Context, email, code string) error {
	s.log.Info().
		Str("email", email).
		Str("code", code).
		Msg("OTP dispatched")
	return nil
}

var _ ports.EmailSender = (*LogSender)(nil)
