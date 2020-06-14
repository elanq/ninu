package main

import (
	"fmt"

	"github.com/elanq/ninu"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	ninu.NewTelegramBot()
	ninu.TelegramBot.Handle("/hello", func(m *tb.Message) {
		ninu.TelegramBot.Send(m.Sender, "halo bro")
	})
	ninu.TelegramBot.Handle("/download", ninu.HandleDownload)

	fmt.Println("Bot is now running")
	ninu.TelegramBot.Start()
	url := "https://twitter.com/LaLigaEN/status/1272107009232396289"
	err := ninu.ProcessURL(url, "thomas_partey.mp4")
	if err != nil {
		panic(err)
	}
}
