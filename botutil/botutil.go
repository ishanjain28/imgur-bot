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
	"fmt"
)

var (
	bot         *tbot.BotAPI
	i           *imgur.Imgur
	redisClient *redis.Client
)

// Initialise global variables in this package
func Init(b *tbot.BotAPI, Imgur *imgur.Imgur, rClient *redis.Client) {
	bot = b
	i = Imgur
	redisClient = rClient
}

// Handle Commands
func HandleCommands(u tbot.Update) {
	cmdArray := strings.Split(u.Message.Text, " ")

	switch cmdArray[0] {
	case "/start":

		msg := tbot.NewMessage(u.Message.Chat.ID, "")
		bot.Send(msg)

	case "/login":
		msgstr := "Open this link in a browser to login:\n"

		msgstr += i.AccessTokenString(strconv.FormatInt(u.Message.Chat.ID, 10) + "-" + u.Message.Chat.UserName)
		msg := tbot.NewMessage(u.Message.Chat.ID, msgstr)
		msg.DisableWebPagePreview = true
		bot.Send(msg)

	case "/stats":

		//Handle the case when user gives a imgur username after /stats command
		if len(cmdArray) > 1 {
			stats, err := i.AccountBase(cmdArray[1], "")
			if err != nil {
				ErrorMessage(u.Message.Chat.ID, err)
				return
			}

			//Create a new common.User with just the username
			user := &common.User{Username: cmdArray[1]}

			UserStatsMessage(u.Message.Chat.ID, stats, nil, nil, user)
			return
		}

		// Fetch user from database
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

		stats, ierr := i.AccountBase(user.Username, "")
		if ierr != nil {
			ErrorMessage(u.Message.Chat.ID, ierr)
			return
		}

		cCount, ierr := i.CommentCount(user.Username, user.AccessToken)
		if ierr != nil {
			ErrorMessage(u.Message.Chat.ID, ierr)
			return
		}

		iCount, ierr := i.ImageCount(user.Username, user.AccessToken)
		if ierr != nil {
			ErrorMessage(u.Message.Chat.ID, ierr)
			return
		}

		UserStatsMessage(u.Message.Chat.ID, stats, cCount, iCount, user)

	case "/logout":
	case "/help":
	default:
		msg := tbot.NewMessage(u.Message.Chat.ID, "Unknown Command, Type /help to get help")
		bot.Send(msg)
	}
}

// Handle Photo uploads
// Fetch User from Database
// Fetch the Albums in that user's Imgur account
// Show a keyboard to user so he/she can select a album or select none at all
// Fetch URL of file
// Upload Image
func HandlePhoto(u tbot.Update) {

	photoSlice := u.Message.Photo

	bestPhoto := (*photoSlice)[2]

	user, err := fetchUser(u.Message.Chat.ID)
	if err != nil {
		if err == redis.Nil {
			UserNotLoggedIn(u.Message.Chat.ID)
			return
		}

		msg := tbot.NewMessage(u.Message.Chat.ID, "Error Occurred, Please retry")
		bot.Send(msg)
		log.Warn.Println("error in fetching user", err.Error())
		return
	}

	albums, ierr := i.Albums(user.Username, user.AccessToken)
	if ierr != nil {
		ErrorMessage(u.Message.Chat.ID, ierr)
		return
	}

	if len(albums.Data) > 0 {
		msg := tbot.NewMessage(u.Message.Chat.ID, "Select an Album")

		var rows [][]tbot.InlineKeyboardButton

		noalbumbtn := tbot.NewInlineKeyboardButtonData("<No Album>", bestPhoto.FileID)
		createbtn := tbot.NewInlineKeyboardButtonData("<Create Album>", bestPhoto.FileID)

		row := [][]tbot.InlineKeyboardButton{{noalbumbtn, createbtn}}

		rows = append(rows, row...)

		for i := 0; i < len(albums.Data); i++ {

			button := tbot.NewInlineKeyboardButtonData(albums.Data[i].Title, bestPhoto.FileID)

			row := [][]tbot.InlineKeyboardButton{{button}}

			rows = append(rows, row...)
		}

		msg.ReplyMarkup = tbot.InlineKeyboardMarkup{InlineKeyboard: rows}

		bot.Send(msg)
	} else {
		// User has no albums, Just upload the image,
		imgUrl, err := bot.GetFileDirectURL(bestPhoto.FileID)
		if err != nil {
			msg := tbot.NewMessage(u.Message.Chat.ID, "error in uploading image, Please retry")
			bot.Send(msg)
			log.Warn.Println("Error in getting file url", err.Error())
		}

		resp, ierr := i.UploadImage(imgUrl, "",
			user.AccessToken)
		if ierr != nil {
			ErrorMessage(u.Message.Chat.ID, ierr)
			return
		}

		msgstr := "Image Uploaded\n"
		msgstr += "URL: " + resp.Data.Link

		msg := tbot.NewMessage(u.Message.Chat.ID, msgstr)
		msg.ReplyToMessageID = u.Message.MessageID
		msg.DisableWebPagePreview = true
		bot.Send(msg)
	}
}

func HandleCallbackQuery(u tbot.Update) {
	user, err := fetchUser(u.CallbackQuery.Message.Chat.ID)
	if err != nil {
		if err == redis.Nil {
			UserNotLoggedIn(u.CallbackQuery.Message.Chat.ID)
			return
		}

		msg := tbot.NewMessage(u.CallbackQuery.Message.Chat.ID, "Error Occurred, Please retry")
		bot.Send(msg)
		log.Warn.Println("error in fetching user", err.Error())
		return
	}

	fileID := u.CallbackQuery.Data

	imgUrl, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		msg := tbot.NewMessage(u.CallbackQuery.Message.Chat.ID, "error in uploading image, Please retry")
		bot.Send(msg)
		log.Warn.Println("Error in getting file url", err.Error())
	}

	resp, ierr := i.UploadImage(imgUrl, u.CallbackQuery.Message.Text, user.AccessToken)
	if ierr != nil {
		ErrorMessage(u.CallbackQuery.Message.Chat.ID, ierr)
		return
	}

	msgstr := "Image Uploaded\n"
	msgstr += "URL: " + resp.Data.Link

	msg := tbot.NewMessage(u.CallbackQuery.Message.Chat.ID, msgstr)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.DisableWebPagePreview = true

	var emptyKeyboardRows [][]tbot.InlineKeyboardButton

	msg.ReplyMarkup = tbot.InlineKeyboardMarkup{InlineKeyboard: emptyKeyboardRows}

	bot.Send(msg)

	fmt.Println(u.CallbackQuery.Message.Text)

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
