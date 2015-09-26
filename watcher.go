package main

import (
	"github.com/ixchi/foxbot/bot"
	"github.com/syfaro/telegram-bot-api"
	"regexp"
	"strconv"
	"strings"
)

type pluginWatcher struct {
	waitingRules *foxbot.WaitingForText
	waitingChan  *foxbot.WaitingForText
}

func (plugin *pluginWatcher) Name() string {
	return "GroupWelcome Watcher"
}

func (plugin *pluginWatcher) GetCommands() []*foxbot.Command {
	return []*foxbot.Command{
		&foxbot.Command{
			Handler:    plugin.allEvents,
			AlwaysCall: true,
		},
		&foxbot.Command{
			Name:    "Rules",
			Help:    "Sends you a list of channel rules",
			Example: "/rules",
			Command: "rules",
			Handler: plugin.rules,
			Waiting: plugin.waitingChan,
		},
		&foxbot.Command{
			Name:    "Set Rules",
			Help:    "Sets the channel rules",
			Example: "/setrules",
			Command: "setrules",
			Handler: plugin.setRules,
			Waiting: plugin.waitingRules,
		},
	}
}

func (plugin *pluginWatcher) allEvents(handler foxbot.Handler) error {
	if !storage.IsEnabled(handler.Update.Message.Chat.ID) {
		return nil
	}

	message := handler.Update.Message

	storage.Set(message.Chat.ID, "name", message.Chat.Title)
	c := storage.Get(message.Chat.ID)

	if message.GroupChatCreated {
		msg := tgbotapi.NewMessage(message.Chat.ID, c["start"])
		_, err := handler.API.SendMessage(msg)
		if err != nil {
			return err
		}
	}

	if message.NewChatParticipant.ID != 0 {
		text := strings.Replace(c["new"], "USER_NAME", message.NewChatParticipant.String(), -1)
		text = strings.Replace(text, "RULES_LINK", "http://groupwelcomebot.xyz/rules/"+strconv.Itoa(message.Chat.ID), -1)

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err := handler.API.SendMessage(msg)
		if err != nil {
			return err
		}
	}

	if message.LeftChatParticipant.ID != 0 {
		text := strings.Replace(c["left"], "USER_NAME", message.NewChatParticipant.String(), -1)
		text = strings.Replace(text, "RULES_LINK", "http://groupwelcomebot.xyz/rules/"+strconv.Itoa(message.Chat.ID), -1)

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err := handler.API.SendMessage(msg)
		return err
	}

	storage.AddChannelToUser(handler.Update.Message.Chat.Title, handler.Update.Message.Chat.ID, handler.Update.Message.From.ID)

	return nil
}

func (plugin *pluginWatcher) setRules(handler foxbot.Handler) error {
	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, "Please enter the new text for the rules:")
	msg.ReplyToMessageID = handler.Update.Message.MessageID
	msg.ReplyMarkup = &tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}

	plugin.waitingRules = &foxbot.WaitingForText{
		IsWaiting: true,
		ChatID:    handler.Update.Message.Chat.ID,
		UserID:    handler.Update.Message.From.ID,
		AnyInChat: false,
		Handler:   plugin.saveRules,
	}

	_, err := handler.API.SendMessage(msg)
	return err
}

func (plugin *pluginWatcher) saveRules(handler foxbot.Handler) error {
	return storage.Set(handler.Update.Message.Chat.ID, "rules", handler.Update.Message.Text)
}

func (plugin *pluginWatcher) rules(handler foxbot.Handler) error {
	if handler.Update.Message.IsGroup() {
		msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, "Please visit http://groupwelcomebot.xyz/rules/"+strconv.Itoa(handler.Update.Message.Chat.ID))
		msg.ReplyToMessageID = handler.Update.Message.MessageID

		_, err := handler.API.SendMessage(msg)
		return err
	}

	myChannels := storage.GetUserChannels(handler.Update.Message.From.ID)
	var keyboard [][]string

	for k, c := range myChannels {
		keyboard = append(keyboard, []string{c + " (" + k + ")"})
	}

	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, "For which channel would you like to see rules?")
	msg.ReplyMarkup = &tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       true,
		Keyboard:        keyboard,
	}

	plugin.waitingChan = &foxbot.WaitingForText{
		IsWaiting: true,
		ChatID:    handler.Update.Message.Chat.ID,
		UserID:    handler.Update.Message.From.ID,
		AnyInChat: false,
		Handler:   plugin.getRules,
	}

	_, err := handler.API.SendMessage(msg)
	return err
}

func (plugin *pluginWatcher) getRules(handler foxbot.Handler) error {
	msg := tgbotapi.NewMessage(handler.Update.Message.Chat.ID, "")
	msg.ReplyMarkup = &tgbotapi.ReplyKeyboardHide{
		HideKeyboard: true,
		Selective:    true,
	}

	r := regexp.MustCompile(`(-\d+)`)
	ch := r.FindString(handler.Update.Message.Text)

	c := storage.GetByString(ch)

	if val, ok := c["rules"]; ok {
		msg.Text = val
	} else {
		msg.Text = "Huh, we don't seem to have this channel."
	}

	_, err := handler.API.SendMessage(msg)
	return err
}
