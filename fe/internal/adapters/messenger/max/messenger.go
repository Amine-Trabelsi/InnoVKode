package max

import (
	"context"
	"fmt"
	"strings"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog"

	"github.com/escalopa/inno-vkode/internal/domain"
	"github.com/escalopa/inno-vkode/internal/ports"
)

type Messenger struct {
	api *maxbot.Api
	log zerolog.Logger
}

var _ ports.Messenger = (*Messenger)(nil)

func New(api *maxbot.Api, log zerolog.Logger) *Messenger {
	return &Messenger{
		api: api,
		log: log,
	}
}

func (m *Messenger) Start(ctx context.Context, handler func(context.Context, domain.Update) error) error {
	updates := m.api.GetUpdates(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case upd, ok := <-updates:
			if !ok {
				return nil
			}
			dUpdate, ok := m.normalizeUpdate(upd)
			if !ok {
				continue
			}
			if err := handler(ctx, dUpdate); err != nil {
				m.log.Error().Err(err).Msg("bot handler error")
			}
			if cb, ok := upd.(*schemes.MessageCallbackUpdate); ok {
				go m.answerCallback(ctx, cb.Callback.CallbackID)
			}
		}
	}
}

func (m *Messenger) normalizeUpdate(upd schemes.UpdateInterface) (domain.Update, bool) {
	switch u := upd.(type) {
	case *schemes.MessageCreatedUpdate:
		return domain.Update{
			Type:      domain.UpdateTypeMessage,
			ChatID:    u.Message.Recipient.ChatId,
			UserID:    u.Message.Sender.UserId,
			Text:      strings.TrimSpace(u.Message.Body.Text),
			MessageID: u.Message.Body.Mid,
			Raw:       upd,
		}, true
	case *schemes.MessageCallbackUpdate:
		var chatID int64
		var mid string
		if u.Message != nil {
			chatID = u.Message.Recipient.ChatId
			mid = u.Message.Body.Mid
		}
		return domain.Update{
			Type:      domain.UpdateTypeCallback,
			ChatID:    chatID,
			UserID:    u.Callback.User.UserId,
			Payload:   u.Callback.Payload,
			MessageID: mid,
			Raw:       upd,
		}, true
	default:
		return domain.Update{}, false
	}
}

func (m *Messenger) answerCallback(ctx context.Context, callbackID string) {
	if callbackID == "" {
		return
	}
	_, err := m.api.Messages.AnswerOnCallback(ctx, callbackID, &schemes.CallbackAnswer{
		Notification: "✅",
	})
	if err != nil {
		m.log.Warn().Err(err).Msg("failed to answer callback")
	}
}

func (m *Messenger) Send(ctx context.Context, chatID, userID int64, msg domain.OutgoingMessage) error {
	builder := maxbot.NewMessage().SetText(msg.Text)
	if chatID > 0 {
		builder.SetChat(chatID)
	}
	// if userID > 0 {
	// 	builder.SetUser(userID)
	// }
	if chatID == 0 && userID == 0 {
		return fmt.Errorf("no chatID/userID provided for outgoing message")
	}
	if msg.ParseMode != "" {
		builder.SetFormat(msg.ParseMode)
	}
	// if msg.Reset {
	// 	builder.SetReset(true)
	// }
	if kb := m.buildKeyboard(msg.Keyboard); kb != nil {
		builder.AddKeyboard(kb)
	}
	_, err := m.api.Messages.Send(ctx, builder)
	if err != nil && err.Error() != "" {
		m.log.Error().Err(err).Msg("failed to send message")
		return err
	}
	return nil
}

func (m *Messenger) buildKeyboard(kb *domain.Keyboard) *maxbot.Keyboard {
	if kb == nil || len(kb.Rows) == 0 {
		return nil
	}
	builder := m.api.Messages.NewKeyboardBuilder()
	for _, row := range kb.Rows {
		r := builder.AddRow()
		for _, btn := range row {
			intent := m.intentFromStyle(btn.Style)
			switch btn.Kind {
			case domain.ButtonKindLink:
				r.AddLink(btn.Label, intent, btn.URL)
			default:
				payload := btn.Payload
				if payload == "" {
					payload = btn.Label
				}
				r.AddCallback(btn.Label, intent, payload)
			}
		}
	}
	return builder
}

func (m *Messenger) intentFromStyle(style domain.ButtonStyle) schemes.Intent {
	switch style {
	case domain.ButtonStylePrimary:
		return schemes.POSITIVE
	case domain.ButtonStyleDanger:
		return schemes.NEGATIVE
	default:
		return schemes.DEFAULT
	}
}
