package main

import (
	"fmt"

	"github.com/elanq/ninu"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	ninu.InitRedis()
	ninu.InitCredential()
	ninu.NewTelegramBot()
	ninu.TelegramBot.Handle("/hello", func(m *tb.Message) {
		ninu.TelegramBot.Send(m.Sender, "halo bro")
	})
	ninu.TelegramBot.Handle("/login", ninu.HandleLogin)
	ninu.TelegramBot.Handle("/authorize", ninu.HandleAuthorize)
	ninu.TelegramBot.Handle("/test", ninu.HandleTest)
	ninu.TelegramBot.Handle("/add", ninu.HandleAdd)
	fmt.Println("Bot is now running")

	ninu.TelegramBot.Start()
}
