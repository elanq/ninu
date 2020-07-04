package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elanq/ninu"
	tb "gopkg.in/tucnak/telebot.v2"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

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

	port := os.Getenv("PORT")
	http.HandleFunc("/healthz", handleHealthz)
	http.ListenAndServe(":"+port, nil)

	ninu.TelegramBot.Start()
}
