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
	UserDelete         = "/userdelete"
	GetUsers           = "/getusers"
	GetUser            = "/getuser"
	StatusSubscription = "/status"
	Subscribe          = "/subscribe"
	DeleteSubscription = "/deletesub"
	GetDevices         = "/getdevices"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	// case WGVpnCmd:
	// 	return p.getConfig(ctx, chatID, username)
	case GetUsers:
		return p.getUsers(ctx, chatID)
	case GetUser:
		return p.getUser(ctx, chatID)
	case UserDelete:
		return p.deleteUser(ctx, chatID)
	case StatusSubscription:
		return p.checkSubscription(ctx, chatID)
	case Subscribe:
		return p.subscribe(ctx, chatID)
	case DeleteSubscription:
		return p.deleteSubscription(ctx, chatID)
	case GetDevices:
		return p.getDevices(ctx, chatID)
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
		return p.tg.SendMessage(ctx, chatID, msgErrorUnknown)
	}
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

// func (p *Processor) getStatusSubscription(ctx context.Context, chatID int, username string) (err error) {
// 	msg, err := p.service.StatusSubscribtion(ctx, username, chatID)
// 	if err != nil {
// 		return p.tg.SendMessage(ctx, chatID, msgErrorGetStatus)
// 	}
// 	return p.tg.SendMessage(ctx, chatID, msg)
// }

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

func (p *Processor) getUser(ctx context.Context, chatID int) error {
	user, err := p.service.GetUser(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, user.Username)
}

func (p *Processor) checkSubscription(ctx context.Context, chatID int) error {
	msg, err := p.service.CheckSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}

	return p.tg.SendMessage(ctx, chatID, msg)
}

func (p *Processor) subscribe(ctx context.Context, chatID int) error {
	err := p.service.ActivateSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, msgSubscribe)
}

func (p *Processor) deleteSubscription(ctx context.Context, chatID int) error {
	err := p.service.DeactivateSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, "подписка удалена")
}

func (p *Processor) getDevices(ctx context.Context, chatID int) error {
	devices, err := p.service.GetDevices(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, err.Error())
	}
	return p.tg.SendMessage(ctx, chatID, devices)
}
