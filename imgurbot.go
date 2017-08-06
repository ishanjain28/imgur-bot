package main

import (
	"os"
	"github.com/ishanjain28/imgur-bot/log"
	tbot    "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"github.com/ishanjain28/imgur-bot/botutil"
	"github.com/ishanjain28/imgur-bot/imgur"
	"github.com/go-redis/redis"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	"github.com/ishanjain28/imgur-bot/common"
)

var (
	GO_ENV              = ""
	PORT                = ""
	TOKEN               = ""
	HOST                = ""
	IMGUR_CLIENT_ID     = ""
	IMGUR_CLIENT_SECRET = ""
	REDIS_URL           = ""
	redisClient         *redis.Client
)

func init() {
	//PORT on which HTTP Server is started
	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Error.Fatalln("$PORT not set")
	}

	//Environment in which application is running
	//In different environments, we should use different methods
	GO_ENV = os.Getenv("GO_ENV")
	if GO_ENV == "" {
		log.Warn.Println("$GO_ENV not set, Setting it to \"development\"")
		GO_ENV = "development"
		HOST = "localhost:" + PORT
	}

	//Telegram Bot Token
	TOKEN = os.Getenv("TOKEN")
	if TOKEN == "" {
		log.Error.Fatalln("Telegram $TOKEN not set")
	}

	IMGUR_CLIENT_ID = os.Getenv("IMGUR_CLIENT_ID")
	if IMGUR_CLIENT_ID == "" {
		log.Error.Fatalln("$IMGUR_CLIENT_ID not set")
	}

	IMGUR_CLIENT_SECRET = os.Getenv("IMGUR_CLIENT_SECRET")
	if IMGUR_CLIENT_SECRET == "" {
		log.Error.Fatalln("$IMGUR_CLIENT_SECRET not set")
	}

	REDIS_URL = os.Getenv("REDISTOGO_URL")
	if REDIS_URL == "" {
		log.Error.Fatalln("$REDISTOGO_URL not set")
	}

	HOST = os.Getenv("HOST")
	if HOST == "" {
		log.Error.Fatalln("$HOST not set")
	}
}

func main() {

	var err error
	//Parse connection string and connect to Database
	redisOpt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Error.Fatalln("error in parsing $REDISTOGO_URL")
	}
	redisClient = redis.NewClient(redisOpt)
	err = redisClient.Ping().Err()
	if err != nil {
		log.Error.Fatalln("error in connecting to DB", err.Error())
	}

	//Initialise imgur
	i, err := imgur.Init(imgur.Config{
		UseFreeAPI:   true,
		ClientID:     IMGUR_CLIENT_ID,
		ClientSecret: IMGUR_CLIENT_SECRET},
	)
	if err != nil {
		log.Error.Fatalln("error in initialising imgur", err.Error())
	}

	//Set OAuth Endpoint
	i.SetOAuthEndpoint("/imgur_oauth", catchImgurOAuthResponse)

	bot, err := tbot.NewBotAPI(TOKEN)
	if err != nil {
		log.Error.Fatalln("error in starting bot", err.Error())
	}

	if GO_ENV == "development" {
		bot.Debug = true
	}
	log.Info.Printf("Authorized on account %s(@%s)\n", bot.Self.FirstName, bot.Self.UserName)

	//Initalise botutil package
	botutil.Init(bot, i, redisClient)

	//Redirect users visting home page to telegram
	http.HandleFunc("/", redirectToTelegram)

	go func() {
		err := http.ListenAndServe(":"+PORT, nil)
		if err != nil {
			log.Error.Fatalln("error in http server", err.Error())
		}
	}()

	//Fetch updates from telegram
	updates := fetchUpdates(bot)

	for update := range updates {

		if update.Message == nil && update.InlineQuery == nil && update.CallbackQuery == nil && update.EditedMessage == nil {
			continue
		}

		handleUpdates(update)
	}
}

func handleUpdates(u tbot.Update) {

	if u.CallbackQuery != nil {
		botutil.HandleCallbackQuery(u)
	}

	if u.Message != nil && u.Message.IsCommand() {
		botutil.HandleCommands(u)
		return
	}

	if u.Message != nil && u.Message.Photo != nil {
		botutil.HandlePhoto(u)
		return
	}

}

func fetchUpdates(bot *tbot.BotAPI) tbot.UpdatesChannel {

	if GO_ENV == "development" {
		//When application is in development mode,
		//Use Polling to fetch updates

		//	Remove any existing webhook
		bot.RemoveWebhook()

		log.Info.Println("Using Polling Method for new updates")
		u := tbot.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Error.Println("error in fetching updates", err.Error())
		}

		return updates
	} else {
		//Use Webhook to receive updates

		//	Remove existing webhook
		bot.RemoveWebhook()

		webhookAddr := "https://" + HOST + "/chinguimgurbot/" + bot.Token
		_, err := bot.SetWebhook(tbot.NewWebhook(webhookAddr))

		if err != nil {
			log.Error.Fatalln("error in setting webhook", webhookAddr, err.Error())
		}

		updates := bot.ListenForWebhook("/chinguimgurbot/" + bot.Token)

		return updates
	}

	return nil
}

func catchImgurOAuthResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	state := r.Form.Get("state")
	expiresin := r.Form.Get("expires_in")
	refToken := r.Form.Get("refresh_token")
	accToken := r.Form.Get("access_token")
	username := r.Form.Get("account_username")
	tusername := strings.Split(state, "-")[1]
	tchatid := strings.Split(state, "-")[0]
	expiresinInt, _ := strconv.Atoi(expiresin)

	serialized, _ := json.Marshal(&common.User{
		TChatID:      tchatid,
		ExpiresIn:    expiresin,
		RefreshToken: refToken,
		AccessToken:  accToken,
		Username:     username,
		TUsername:    tusername,
	})

	res, err := redisClient.Set(tchatid, serialized, time.Duration(expiresinInt)*time.Second).Result()
	if err != nil {
		log.Warn.Println("error in storing login information", err.Error())
	}
	if res == "OK" {
		log.Info.Printf("%s logged in with imgur account %s", tusername, username)
	}
}

func redirectToTelegram(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://t.me/chinguimgurbot", http.StatusTemporaryRedirect)
}
