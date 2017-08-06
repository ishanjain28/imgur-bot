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

// Initialise global variables in this package
func Init(b *tbot.BotAPI, Imgur *imgur.Imgur, rClient *redis.Client) {
	bot = b
	i = Imgur
	redisClient = rClient
}

// Handle Commands
func HandleCommands(u tbot.Update) {
	cmdArray := strings.Split(u.Message.Text, " ")
	chatID := u.Message.Chat.ID

	switch cmdArray[0] {
	case "/start":

		msg := tbot.NewMessage(chatID, "@chinguimgurbot helps in uploading images to imgur.com and in checking stats.\n")
		bot.Send(msg)

	case "/login":

		msgstr := "Open this link in a browser to login:\n"

		msgstr += i.AccessTokenString(strconv.FormatInt(chatID, 10) + "-" + u.Message.Chat.UserName)
		msg := tbot.NewMessage(chatID, msgstr)
		msg.DisableWebPagePreview = true
		bot.Send(msg)

	case "/stats":

		//Handle the case when user gives a imgur username after /stats command
		if len(cmdArray) > 1 {
			stats, err := i.AccountBase(cmdArray[1], "")
			if err != nil {
				ErrorMessage(chatID, err)
				return
			}

			//Create a new common.User with just the username
			user := &common.User{Username: cmdArray[1]}

			UserStatsMessage(chatID, stats, nil, nil, user)
			return
		}

		// Fetch user from database
		user, err := fetchUser(chatID)
		if err != nil {
			if err == redis.Nil {
				UserNotLoggedIn(chatID)
				return
			}

			msg := tbot.NewMessage(chatID, "error in fetching user "+err.Error())
			bot.Send(msg)
			log.Warn.Println("error in fetching user", err.Error())
			return
		}

		stats, ierr := i.AccountBase(user.Username, "")
		if ierr != nil {
			ErrorMessage(chatID, ierr)
			return
		}

		cCount, ierr := i.CommentCount(user.Username, user.AccessToken)
		if ierr != nil {
			ErrorMessage(chatID, ierr)
			return
		}

		iCount, ierr := i.ImageCount(user.Username, user.AccessToken)
		if ierr != nil {
			ErrorMessage(chatID, ierr)
			return
		}

		UserStatsMessage(chatID, stats, cCount, iCount, user)

	case "/logout":

		err := deleteUser(chatID)
		if err != nil {
			msg := tbot.NewMessage(chatID, "Error in signing out, Please retry")
			bot.Send(msg)
			return
		}

		msg := tbot.NewMessage(chatID, "Successfully logged out")
		bot.Send(msg)

	case "/help":

		msgstr := "@chinguimgurbot helps in uploading images to imgur.com and in checking stats.\n"
		msgstr += "To get started, You need to /login\n"
		msgstr += "After logging in, You can upload images by sending it an image or check stats using /stats command"
		msgstr += "You can also check someone else's stats using /stats <A imgur username>"

		msg := tbot.NewMessage(chatID, msgstr)

		bot.Send(msg)

	default:
		msg := tbot.NewMessage(chatID, "Unknown Command, Type /help to get help")
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

		noalbumbtn := tbot.NewInlineKeyboardButtonData("<No Album>", "-1\\"+bestPhoto.FileID)

		row := [][]tbot.InlineKeyboardButton{{noalbumbtn}}

		rows = append(rows, row...)

		for i := 0; i < len(albums.Data); i++ {

			// skip albums with length > 64
			if len(albums.Data[i].Title) > 64 {
				continue
			}

			button := tbot.NewInlineKeyboardButtonData(albums.Data[i].Title, strconv.Itoa(i)+"\\"+bestPhoto.FileID)

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

	chatID := u.CallbackQuery.Message.Chat.ID
	messageID := u.CallbackQuery.Message.MessageID
	datas := strings.Split(u.CallbackQuery.Data, "\\")
	fileID := datas[1]
	albumIndex, err := strconv.Atoi(datas[0])
	if err != nil {
		msg := tbot.NewMessage(chatID, "Error in parsing album index")
		bot.Send(msg)
		log.Warn.Println("Error in fetching album index", err.Error())
		return

	}

	user, err := fetchUser(chatID)
	if err != nil {
		if err == redis.Nil {
			UserNotLoggedIn(chatID)
			return
		}

		msg := tbot.NewMessage(chatID, "Error Occurred, Please retry")
		bot.Send(msg)
		log.Warn.Println("error in fetching user", err.Error())
		return
	}

	var albums *imgur.Album
	var ierr *imgur.IError
	var resp *imgur.Image
	if albumIndex >= 0 {
		albums, ierr = i.Albums(user.Username, user.AccessToken)
		if err != nil {
			msg := tbot.NewMessage(chatID, "Error: "+ierr.String())
			bot.Send(msg)
			log.Warn.Println("Error in getting file url", err.Error())
		}
	}

	imgUrl, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		msg := tbot.NewMessage(chatID, "error in uploading image, Please retry")
		bot.Send(msg)
		log.Warn.Println("Error in getting file url", err.Error())
	}
	msgstr := "Image Uploaded\n"

	if albumIndex >= 0 {
		msgstr += "Album: " + albums.Data[albumIndex].Title + "\n"

		resp, ierr = i.UploadImage(imgUrl, albums.Data[albumIndex].ID, user.AccessToken)
		if ierr != nil {
			ErrorMessage(chatID, ierr)
			return
		}
	} else {
		resp, ierr = i.UploadImage(imgUrl, "", user.AccessToken)
		if ierr != nil {
			ErrorMessage(chatID, ierr)
			return
		}
	}

	msgstr += "URL: " + resp.Data.Link

	// Delete the inline keyboard
	bot.DeleteMessage(tbot.DeleteMessageConfig{MessageID: messageID, ChatID: chatID})

	msg := tbot.NewMessage(chatID, msgstr)
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

func deleteUser(chatid int64) error {
	chatidStr := strconv.FormatInt(chatid, 10)
	_, err := redisClient.Del(chatidStr).Result()
	if err != nil {
		return err
	}

	return nil
}
