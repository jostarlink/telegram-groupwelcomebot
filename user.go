package main

import (
	"bytes"
	"github.com/ixchi/foxbot/bot"
	"github.com/syfaro/telegram-bot-api"
)

type pluginUser struct {
}

func (plugin *pluginUser) Name() string {
	return "GroupWelcome User Controls"
}

func (plugin *pluginUser) GetCommands() []*foxbot.Command {
	return []*foxbot.Command{
		&foxbot.Command{
			Name:    "My Channels",
			Help:    "Displays my channels",
			Example: "/mychannels",
			Command: "mychannels",
			Handler: plugin.myChannels,
		},
		&foxbot.Command{
			Name:    "Clear Channels",
			Help:    "Clears the list of channels you are currently in",
			Example: "/clearchans",
			Command: "clearchans",
			Handler: plugin.myChannels,
		},
	}
}

func (plugin *pluginUser) myChannels(handler foxbot.Handler) error {
	if handler.Update.Message.IsGroup() {
		msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, "Please only run this command in a direct message")
		msg.ReplyToMessageID = handler.Update.Message.MessageID

		_, err := handler.API.SendMessage(msg)
		return err
	}

	myChannels := storage.GetUserChannels(handler.Update.Message.From.ID)

	b := &bytes.Buffer{}

	for _, c := range myChannels {
		b.WriteString(c)
		b.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, b.String())
	_, err := handler.API.SendMessage(msg)

	return err
}
