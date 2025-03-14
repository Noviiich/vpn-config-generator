package telegram

import (
	"context"
	"log"
	"strings"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	OpenVpnCmd = "/openvpn"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case OpenVpnCmd:
		return p.CreateConfig(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) CreateConfig(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't create config", err) }()

	config, err := p.service.Create(ctx, username)
	if err != nil {
		return err
	}

	if err := p.tg.SendMessage(ctx, chatID, config); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
