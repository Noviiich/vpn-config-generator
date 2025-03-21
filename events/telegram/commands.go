package telegram

import (
	"context"
	"log"
	"strings"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

const (
	HelpCmd   = "/help"
	StartCmd  = "/start"
	WGVpnCmd  = "/wireguard"
	VpnStatus = "/status"
	VpnSub    = "/subscribe"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case WGVpnCmd:
		return p.createConfig(ctx, chatID, username)
	case VpnStatus:
		return p.statusSubscription(ctx, chatID, username)
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

func (p *Processor) createConfig(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't create config", err) }()

	configText, err := p.service.Create(ctx, username, chatID)
	if err != nil {
		p.tg.SendMessage(ctx, chatID, err.Error())
		return err
	}

	if err := p.tg.SendDocument(ctx, chatID, configText, "Wireguard.conf"); err != nil {
		p.tg.SendMessage(ctx, chatID, err.Error())
		return err
	}

	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int, username string) error {
	err := p.service.CreateUser(ctx, username, chatID)
	if err != nil {
		p.tg.SendMessage(ctx, chatID, msgErrorCreateUser)
	}
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func (p *Processor) statusSubscription(ctx context.Context, chatID int, username string) (err error) {
	msg, err := p.service.StatusSubscribtion(ctx, username, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorGetStatus)
	}
	return p.tg.SendMessage(ctx, chatID, msg)
}

func (p *Processor) subscribe(ctx context.Context, chatID int) error {
	err := p.service.UpdateSubscription(ctx, chatID)
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, msgErrorSubscribe)
	}
	return p.tg.SendMessage(ctx, chatID, msgSubscribe)
}
