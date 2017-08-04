package botutil

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/imgur-bot/imgur"
	"strconv"
	"github.com/go-redis/redis"
	"github.com/ishanjain28/imgur-bot/common"
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/log"
	"strings"
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
	cmdArray := strings.Split(u.Message.Text, " ")

	switch cmdArray[0] {
	case "/start":

		msg := tbot.NewMessage(u.Message.Chat.ID, "")
		bot.Send(msg)

	case "/login":
		msg := tbot.NewMessage(u.Message.Chat.ID, i.AccessTokenString(strconv.FormatInt(u.Message.Chat.ID, 10)+"-"+u.Message.Chat.UserName))
		bot.Send(msg)

	case "/stats":

		//Handle the case when user gives a imgur username after /stats command

		if len(cmdArray) > 1 {
			stats, err := i.AccountBase(cmdArray[1], "")
			if err != nil {
				ErrorResponse(u.Message.Chat.ID, err)
				return
			}

			//Create a new common.User with just the username
			user := &common.User{Username: cmdArray[1]}

			UserStatsMessage(u.Message.Chat.ID, stats, user)
			return
		}

		user, err := fetchUser(u.Message.Chat.ID)
		if err != nil {

			if err == redis.Nil {
				UserNotLoggedIn(u.Message.Chat.ID)
				return
			}

			msg := tbot.NewMessage(u.Message.Chat.ID, "error in fetching user "+err.Error())
			bot.Send(msg)
			log.Warn.Println("error in fetching user", err.Error())
			return
		}
		stats, err := i.AccountBase(user.Username, "")
		if err != nil {
			ErrorResponse(u.Message.Chat.ID, err)
			return
		}

		UserStatsMessage(u.Message.Chat.ID, stats, user)

	case "/logout":
	case "/help":
	default:
		msg := tbot.NewMessage(u.Message.Chat.ID, "Unknown Command, Type /help to get help")
		bot.Send(msg)
	}
}

func HandlePhoto(u tbot.Update) {

	photoSlice := u.Message.Photo

	bestPhoto := (*photoSlice)[2]

	user, err := fetchUser(u.Message.Chat.ID)
	if err != nil {

		if err == redis.Nil {
			UserNotLoggedIn(u.Message.Chat.ID)
			return
		}

		msg := tbot.NewMessage(u.Message.Chat.ID, "error in fetching user "+err.Error())
		bot.Send(msg)
		log.Warn.Println("error in fetching user", err.Error())
		return
	}

	imgUrl, err := bot.GetFileDirectURL(bestPhoto.FileID)
	if err != nil {
		msg := tbot.NewMessage(u.Message.Chat.ID, "error in uploading image, Please retry")
		bot.Send(msg)
	}

	resp, err := i.UploadImage(imgUrl, user)
	if err != nil {
		ErrorResponse(u.Message.Chat.ID, err)
		return
	}

	msgstr := "Image Uploaded\n"
	msgstr += "URL: " + resp.Data.Link

	msg := tbot.NewMessage(u.Message.Chat.ID, msgstr)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.DisableWebPagePreview = true
	bot.Send(msg)
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
