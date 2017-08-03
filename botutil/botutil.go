package botutil

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/imgur-bot/imgur"
)

func HandleCommands(bot *tbot.BotAPI, i *imgur.Imgur, u tbot.Update) {

	switch u.Message.Text {
	case "/start":

		msg := tbot.NewMessage(u.Message.Chat.ID, "")
		bot.Send(msg)

	case "/login":
	case "/logout":
	case "/help":
	default:
		msg := tbot.NewMessage(u.Message.Chat.ID, i.AccessTokenString(""))
		bot.Send(msg)
	}
}
