package main

import (
	"github.com/ixchi/foxbot/bot"
	"github.com/syfaro/telegram-bot-api"
	"io/ioutil"
)

type pluginHelp struct {
}

func (plugin *pluginHelp) Name() string {
	return "GroupWelcome Help"
}

func (plugin *pluginHelp) GetCommands() []*foxbot.Command {
	return []*foxbot.Command{
		&foxbot.Command{
			Name:    "Bot Help",
			Help:    "Help for using this bot",
			Example: "/help",
			Command: "help",
			Handler: plugin.botHelp,
		},
	}
}

func (plugin *pluginHelp) botHelp(handler foxbot.Handler) error {
	b, err := ioutil.ReadFile("bothelp.txt")
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, string(b))
	msg.ReplyToMessageID = handler.Update.Message.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err = handler.API.SendMessage(msg)
	return err
}
