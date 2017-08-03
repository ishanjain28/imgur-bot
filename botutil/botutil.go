package botutil

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/imgur-bot/imgur"
	"strconv"
	"github.com/go-redis/redis"
	"github.com/ishanjain28/imgur-bot/common"
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/log"
)

var (
	bot         *tbot.BotAPI
	i           *imgur.Imgur
	redisClient *redis.Client
)

func Init(b *tbot.BotAPI, Imgur *imgur.Imgur, rClient *redis.Client) {
	bot = b
	i = Imgur
	redisClient = rClient
}

func HandleCommands(u tbot.Update) {

	switch u.Message.Text {
	case "/start":

		msg := tbot.NewMessage(u.Message.Chat.ID, "")
		bot.Send(msg)

	case "/login":
		msg := tbot.NewMessage(u.Message.Chat.ID, i.AccessTokenString(strconv.FormatInt(u.Message.Chat.ID, 10)+"-"+u.Message.Chat.UserName))
		bot.Send(msg)

	case "/stats":

		user, err := fetchUser(u.Message.Chat.ID)
		if err != nil {

			if err == redis.Nil {
				//todo:COMPLETE THIS
			}

			msg := tbot.NewMessage(u.Message.Chat.ID, "Error in fetching user "+err.Error())
			bot.Send(msg)
			log.Warn.Println("Error in fetching user", err.Error())
			return
		}
		stats, err := i.AccountBase(user.Username, "")
		if err != nil {
			msg := tbot.NewMessage(u.Message.Chat.ID, "Error in fetching stats "+err.Error())
			bot.Send(msg)
			log.Warn.Println("Error in fetching stats", err.Error())
			return
		}

		msg := tbot.NewMessage(u.Message.Chat.ID, stats.Data.URL)
		bot.Send(msg)

	case "/logout":
	case "/help":
	default:
		msg := tbot.NewMessage(u.Message.Chat.ID, "Unknown Command, Type /help to get help")
		bot.Send(msg)
	}
}

func fetchUser(chatid int64) (*common.User, error) {
	chatidStr := strconv.FormatInt(chatid, 10)
	ustr, err := redisClient.Get(chatidStr).Result()
	if err != nil {
		return nil, err
	}
	a := &common.User{}

	err = json.Unmarshal([]byte(ustr), a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
