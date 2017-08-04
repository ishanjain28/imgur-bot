package botutil

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/imgur-bot/imgur"
	"github.com/ishanjain28/imgur-bot/common"
	"github.com/ishanjain28/imgur-bot/log"
	"strconv"
	"time"
	"fmt"
)

func UserNotLoggedIn(cid int64) {
	msg := tbot.NewMessage(cid, "You are not logged in, Type /login to login")
	bot.Send(msg)
}

func UserStatsMessage(cid int64, stats *imgur.AccountBase, cCount *imgur.Basic, iCount *imgur.Basic, user *common.User) {

	if cCount != nil && iCount != nil {
		fmt.Println(cCount.Data, iCount.Data)
	}
	msgstr := "*Username*: " + user.Username + "\n"

	msgstr += "*Reputation*: " + strconv.Itoa(stats.Data.Reputation) + " (_" + stats.Data.ReputationName + " _)\n"

	msgstr += "*Profile Link*: http://imgur.com/user/" + stats.Data.URL + "\n"

	if stats.Data.Avatar != nil {
		msgstr += "*Avatar*: " + stats.Data.Avatar.(string) + "\n"
	}
	createdOn, err := time.ParseDuration(strconv.FormatInt(time.Now().Unix()-int64(stats.Data.Created), 10) + "s")
	if err != nil {
		msgstr += "*Account created on* " + strconv.Itoa(stats.Data.Created) + "\n"

		log.Warn.Println(err.Error())
	} else {
		msgstr += "*Account was created " + strconv.FormatFloat(createdOn.Hours(), 'f', 1, 64) + "* Hours Ago\n"
	}

	if stats.Data.UserFollow.Status {
		msgstr += "You *follow* this user on imgur.com\n"
	} else {
		msgstr += "You *do not follow* this user on imgur.com\n"
	}

	if stats.Data.Bio != nil {
		msgstr += "*Bio*:\n" + stats.Data.Bio.(string)
	}

	msg := tbot.NewMessage(cid, msgstr)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}

func ErrorResponse(cid int64, err error) {

	msg := tbot.NewMessage(cid, err.Error())
	bot.Send(msg)
	log.Warn.Println("Some error Occurred, Please retry", err.Error())
}
