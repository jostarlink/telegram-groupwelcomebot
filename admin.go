package main

import (
	"bytes"
	"github.com/ixchi/foxbot/bot"
	"github.com/syfaro/telegram-bot-api"
	"strings"
)

type pluginAdmin struct {
}

func (plugin *pluginAdmin) Name() string {
	return "GroupWelcome Admin"
}

func (plugin *pluginAdmin) start(handler foxbot.Handler) error {
	message := handler.Update.Message

	err := storage.Start(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Initialized bot for this channel!\nPlease type /enable to start sending messages.\nFor configuration help, please type /bothelp.")
	msg.ReplyToMessageID = message.MessageID

	_, err = handler.API.SendMessage(msg)
	return err
}

func (plugin *pluginAdmin) enable(handler foxbot.Handler) error {
	return storage.Enable(handler.Update.Message.Chat.ID)
}

func (plugin *pluginAdmin) disable(handler foxbot.Handler) error {
	return storage.Disable(handler.Update.Message.Chat.ID)
}

func (plugin *pluginAdmin) set(handle foxbot.Handler) error {
	return storage.Set(handle.Update.Message.Chat.ID, handle.Args[0], strings.Join(handle.Args[1:len(handle.Args)], " "))
}

func (plugin *pluginAdmin) save(handler foxbot.Handler) error {
	return storage.Save()
}

func (plugin *pluginAdmin) botFather(handler foxbot.Handler) error {
	text := &bytes.Buffer{}
	for _, plugin := range bot.Plugins {
		for _, cmd := range plugin.GetCommands() {
			if cmd.Command != "" && cmd.Name != "" && cmd.Help != "" {
				text.WriteString(cmd.Command)
				text.WriteString(" - ")
				text.WriteString(cmd.Help)
				text.WriteString("\n")
			}
		}
	}

	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, text.String())
	_, err := handler.API.SendMessage(msg)

	return err
}

func (plugin *pluginAdmin) GetCommands() []*foxbot.Command {
	return []*foxbot.Command{
		&foxbot.Command{
			Name:    "Start",
			Help:    "Initializes the bot in the current channel",
			Example: "/start",
			Command: "start",
			Handler: plugin.start,
		},
		&foxbot.Command{
			Name:    "Enable",
			Help:    "Enables the bot in the current channel",
			Example: "/enable",
			Command: "enable",
			Handler: plugin.enable,
		},
		&foxbot.Command{
			Name:    "Stop",
			Help:    "Stops the bot from running in the current channel",
			Example: "/stop",
			Command: "stop",
			Handler: plugin.disable,
		},
		&foxbot.Command{
			Name:    "Save",
			Help:    "Forces a config save",
			Example: "/save",
			Command: "save",
			Handler: plugin.save,
		},
		&foxbot.Command{
			Name:    "Set",
			Help:    "Sets a config option",
			Example: "/set new Hello, new user called USER_NAME!",
			Command: "set",
			Handler: plugin.set,
		},
		&foxbot.Command{
			Name:    "BotFather Help",
			Help:    "Displays all commands, formatted for BotFather",
			Example: "/botfather",
			Command: "botfather",
			Handler: plugin.botFather,
		},
	}
}
