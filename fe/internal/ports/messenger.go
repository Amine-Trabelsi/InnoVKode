package ports

import (
	"context"

	"github.com/escalopa/inno-vkode/internal/domain"
)

type Messenger interface {
	Start(ctx context.Context, handler func(context.Context, domain.Update) error) error
	Send(ctx context.Context, chatID, userID int64, msg domain.OutgoingMessage) error
}
