package main

import (
	"os"
	"github.com/ishanjain28/imgur-bot/log"
	tbot    "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"github.com/ishanjain28/imgur-bot/botutil"
	"github.com/ishanjain28/imgur-bot/imgur"
)

var (
	GO_ENV              = ""
	PORT                = ""
	TOKEN               = ""
	HOST                = "chinguimgurbot.herokuapp.com:80"
	IMGUR_CLIENT_ID     = ""
	IMGUR_CLIENT_SECRET = ""
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
}

func main() {

	//Initalise imgur
	i, err := imgur.Init(imgur.Config{
		UseFreeAPI:   true,
		ClientID:     IMGUR_CLIENT_ID,
		ClientSecret: IMGUR_CLIENT_SECRET},
	)
	if err != nil {
		log.Error.Fatalln("Error in initialising imgur", err.Error())
	}

	//Set OAuth Endpoint
	i.SetOAuthEndpoint("/imgur_oauth")

	bot, err := tbot.NewBotAPI(TOKEN)
	if err != nil {
		log.Error.Fatalln("Error in starting bot", err.Error())
	}

	if GO_ENV == "development" {
		//bot.Debug = true
	}
	log.Info.Printf("Authorized on account %s(@%s)\n", bot.Self.FirstName, bot.Self.UserName)

	//Start HTTP Server
	// There are two uses of this server.
	//1. Receive messages from telegram when in production environment(it uses polling in development)
	//2. Serve as the webhook in imgur's oauth authentication of user

	go func() {
		err := http.ListenAndServe(":"+PORT, nil)
		if err != nil {
			log.Error.Fatalln("Error in http server", err.Error())
		}
	}()

	//Fetch updates from telegram
	updates := fetchUpdates(bot)

	for update := range updates {

		if update.Message == nil && update.InlineQuery == nil && update.CallbackQuery == nil && update.EditedMessage == nil {
			continue
		}

		handleUpdates(bot, i, update)
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
			log.Error.Println("Error in fetching updates", err.Error())
		}

		return updates
	} else {
		//Use Webhook to receive updates

		//	Remove existing webhook
		bot.RemoveWebhook()

		webhookAddr := "https://" + HOST + "/chinguimgurbot/" + bot.Token
		_, err := bot.SetWebhook(tbot.NewWebhook(webhookAddr))

		if err != nil {
			log.Error.Fatalln("Error in setting webhook", webhookAddr, err.Error())
		}

		//updates := bot.ListenForWebhook("/chinguimgurbot/" + bot.Token)

		//	TODO Complete this.

	}

	return nil
}

func handleUpdates(bot *tbot.BotAPI, i *imgur.Imgur, u tbot.Update) {

	if u.Message.IsCommand() {
		botutil.HandleCommands(bot, i, u)
	}

}
