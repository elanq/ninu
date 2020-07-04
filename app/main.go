package main

import (
	"log"
	"net/http"
	"os"

	"github.com/elanq/ninu"
	tb "gopkg.in/tucnak/telebot.v2"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	ninu.InitPostgre()
	ninu.InitRedis()
	ninu.InitCredential()
	ninu.NewTelegramBot()

	ninu.TelegramBot.Handle("/hello", func(m *tb.Message) {
		ninu.TelegramBot.Send(m.Sender, "halo bro")
	})
	ninu.TelegramBot.Handle("/login", ninu.HandleLogin)
	ninu.TelegramBot.Handle("/authorize", ninu.HandleAuthorize)
	ninu.TelegramBot.Handle("/add", ninu.HandleAdd)
	ninu.TelegramBot.Handle("/show", ninu.HandleShow)
	log.Println("Bot is now running")
	go ninu.TelegramBot.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/healthz", handleHealthz)
	log.Println("Server is now running at port", port)
	http.ListenAndServe(":"+port, nil)
}
