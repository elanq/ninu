package ninu

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var TelegramBot *tb.Bot

func NewTelegramBot() {
	botSettings := tb.Settings{
		Token:  "secret",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tb.NewBot(botSettings)
	if err != nil {
		panic(err)
	}
	TelegramBot = b
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
