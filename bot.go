package ninu

import (
	"fmt"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	TelegramBot *tb.Bot
)

func NewTelegramBot() {
	token := os.Getenv("TELEGRAM_API_TOKEN")
	botSettings := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tb.NewBot(botSettings)
	if err != nil {
		panic(err)
	}
	TelegramBot = b
}

func HandleLogin(message *tb.Message) {
	if token, _ := SavedToken(); token != nil {
		TelegramBot.Send(message.Sender, "You are already authorized")
		return
	}

	authURL := AuthURL()
	msg := fmt.Sprintf("go to this URL and copy the authorization code %v", authURL)
	TelegramBot.Send(message.Sender, msg)
}

func HandleAuthorize(message *tb.Message) {
	sender := message.Sender
	authCode := message.Payload
	if err := Authorize(authCode); err != nil {
		TelegramBot.Send(sender, "Error while authorizing code")
		TelegramBot.Send(sender, err.Error())
		return
	}

	TelegramBot.Send(sender, "User authorized")
}

func HandleShow(message *tb.Message) {
	switch message.Payload {
	case "today":
		msg, err := ShowTodayTransaction()
		if err != nil {
			TelegramBot.Send(message.Sender, err.Error())
			return
		}
		TelegramBot.Send(message.Sender, msg)
	}
}

func HandleAdd(message *tb.Message) {
	if err := AddTransaction(message.Payload); err != nil {
		TelegramBot.Send(message.Sender, err.Error())
		return
	}

	TelegramBot.Send(message.Sender, "Transaction added")
}

func HandleDownload(message *tb.Message) {
	sender := message.Sender
	senderID := sender.ID
	messageID := message.ID
	filename := fmt.Sprintf("%d-%d", senderID, messageID)
	url := message.Payload

	TelegramBot.Send(sender, "Please wait, your message is being processed")
	err := ProcessURL(url, filename)
	if err != nil {
		TelegramBot.Send(sender, "There's an error while procesing your message")
		TelegramBot.Send(sender, err)
		return
	}

	videoFile := &tb.Video{File: tb.FromDisk(filename)}
	_, err = TelegramBot.Send(sender, videoFile)
	if err != nil {
		TelegramBot.Send(sender, "error while sending the video to telegram")
	}
}
