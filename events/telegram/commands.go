package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	WGVpnCmd   = "/wireguard"
	VpnStatus  = "/status"
	VpnSub     = "/subscribe"
	UserDelete = "/userdelete"
	GetUsers   = "/getusers"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case WGVpnCmd:
		return p.getConfig(ctx, chatID, username)
	case VpnStatus:
		return p.getStatusSubscription(ctx, chatID, username)
	case GetUsers:
		return p.getUsers(ctx, chatID)
	case UserDelete:
		return p.deleteUser(ctx, chatID)
	case VpnSub:
		return p.subscribe(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID, username)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) getConfig(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't get config", err) }()

	configText, err := p.service.GetConfig(ctx, chatID, username)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorGetConfig)
	}
	if err == nil && configText == "" {
		return p.tg.SendMessage(ctx, chatID, msgNoSubscription)
	}

	if err := p.tg.SendDocument(ctx, chatID, configText, "WG_NOV.conf"); err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorSendDocument)
	}

	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int, username string) error {
	err := p.service.CreateUser(ctx, username, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorCreateUser)
	}
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func (p *Processor) getStatusSubscription(ctx context.Context, chatID int, username string) (err error) {
	msg, err := p.service.StatusSubscribtion(ctx, username, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorGetStatus)
	}
	return p.tg.SendMessage(ctx, chatID, msg)
}

func (p *Processor) subscribe(ctx context.Context, chatID int, username string) error {
	err := p.tg.SendMessage(ctx, chatID, "Ждите подтверждение администратора")
	if err != nil {
		return err
	}
	return p.tg.SendApprovalButtons(ctx, fmt.Sprintf("Запрос на подтверждение : @%s", username), chatID)
}

func (p *Processor) deleteUser(ctx context.Context, chatID int) error {
	err := p.service.DeleteUser(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorDeleteUser)
	}
	return p.tg.SendMessage(ctx, chatID, msgDeleteUser)
}

func (p *Processor) getUsers(ctx context.Context, chatID int) error {
	users, err := p.service.GetUsers(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorGetUsers)
	}
	if users == nil {
		return p.tg.SendMessage(ctx, chatID, msgNoUsers)
	}
	var usernames []string
	for _, user := range users {
		usernames = append(usernames, "@"+user.Username)
	}
	result := strings.Join(usernames, "\n")

	return p.tg.SendMessage(ctx, chatID, result)
}
