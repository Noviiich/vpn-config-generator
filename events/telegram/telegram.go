package telegram

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Noviiich/vpn-config-generator/clients/telegram"
	"github.com/Noviiich/vpn-config-generator/events"
	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/service"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	service *service.VPNService
}

type Meta struct {
	MessageID int
	ChatID    int
	Username  string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, service *service.VPNService) *Processor {
	return &Processor{
		tg:      client,
		service: service,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	case events.CallbackQuery:
		return p.processCallbackQuery(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (p *Processor) processCallbackQuery(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	parts := strings.Split(event.Text, "_")
	if len(parts) != 2 {
		return errors.New("invalid callback data format")
	}

	action := parts[0]
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		return e.Wrap("invalid user ID in callback data", err)
	}

	err = p.tg.DeleteApprovalButtons(ctx, meta.MessageID)
	if err != nil {
		return e.Wrap("can't delete approval buttons", err)
	}

	switch action {
	case "approve":
		if err := p.service.UpdateSubscription(ctx, userID); err != nil {
			return e.Wrap("can't approve subscription", err)
		}
	case "reject":
		return p.tg.NotifyUserSubscriptionRejected(ctx, userID)
	default:
		return errors.New("unknown callback action")
	}

	return p.tg.NotifyUserSubscriptionApproved(ctx, userID)
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			MessageID: upd.Message.MessageID,
			ChatID:    upd.Message.Chat.ID,
			Username:  upd.Message.From.Username,
		}
	} else if updType == events.CallbackQuery {
		res.Meta = Meta{
			MessageID: upd.CallbackQuery.Message.MessageID,
			ChatID:    upd.CallbackQuery.Message.Chat.ID,
			Username:  upd.CallbackQuery.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	} else if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	} else if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}
	return events.Unknown
}
