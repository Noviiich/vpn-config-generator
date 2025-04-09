package telegram

import (
	"context"
	"log"
	"strings"
)

const (
	HelpCmd            = "/help"
	StartCmd           = "/start"
	WGVpnCmd           = "/wireguard"
	VpnStatus          = "/status"
	VpnSub             = "/subscribe"
	UserDelete         = "/userdelete"
	GetUsers           = "/getusers"
	DeleteSubscription = "/deletesub"
	CreateWireguard    = "/wg"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	// case WGVpnCmd:
	// 	return p.getConfig(ctx, chatID, username)
	case VpnStatus:
		return p.getSubscription(ctx, chatID)
	case CreateWireguard:
		return p.createWireguard(ctx, chatID)
	case DeleteSubscription:
		return p.DeleteSubscription(ctx, chatID)
	case GetUsers:
		return p.getUsers(ctx, chatID)
	case UserDelete:
		return p.deleteUser(ctx, chatID)
	case VpnSub:
		return p.subscribe(ctx, chatID)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID, username)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

// func (p *Processor) getConfig(ctx context.Context, chatID int, username string) (err error) {
// 	defer func() { err = e.WrapIfErr("can't do command: can't get config", err) }()

// 	configText, err := p.service.GetConfig(ctx, chatID, username)
// 	if err != nil {
// 		return p.tg.SendMessage(ctx, chatID, msgErrorGetConfig)
// 	}
// 	if err == nil && configText == "" {
// 		return p.tg.SendMessage(ctx, chatID, msgNoSubscription)
// 	}

// 	if err := p.tg.SendDocument(ctx, chatID, configText, "WG_NOV.conf"); err != nil {
// 		return p.tg.SendMessage(ctx, chatID, msgErrorSendDocument)
// 	}

// 	return nil
// }

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

func (p *Processor) getSubscription(ctx context.Context, chatID int) (err error) {
	msg, err := p.service.GetSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msg)
}

func (p *Processor) subscribe(ctx context.Context, chatID int) error {
	err := p.service.CreateSubscription(ctx, chatID, 1)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msgSubscribe)
}

func (p *Processor) deleteUser(ctx context.Context, chatID int) error {
	err := p.service.DeleteUser(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msgDeleteUser)
}

func (p *Processor) getUsers(ctx context.Context, chatID int) error {
	users, err := p.service.GetUsers(ctx)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, users)
}

func (p *Processor) DeleteSubscription(ctx context.Context, chatID int) error {
	err := p.service.DeleteSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msgDeleteSubscription)
}

func (p *Processor) createWireguard(ctx context.Context, chatID int) error {
	err := p.service.CreateAction(ctx, chatID, 1)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msgCrateConfig)
}
