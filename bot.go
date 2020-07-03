package ninu

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/api/sheets/v4"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	TelegramBot  *tb.Bot
	sheetService *sheets.Service
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

func sheetClient() (*sheets.Service, error) {
	if sheetService != nil {
		return sheetService, nil
	}

	return sheets.New(Client())
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

func HandleAdd(message *tb.Message) {
	if Client() == nil {
		TelegramBot.Send(message.Sender, "Client is nil, authorize first")
	}

	client, err := sheetClient()
	if err != nil {
		TelegramBot.Send(message.Sender, err.Error())
	}

	sheetID := os.Getenv("SPREADSHEET_ID")
	sheetRange := "ELANQIST0609_1137757232!I:K"

	in, err := ReadTransaction(message.Payload)
	if err != nil {
		TelegramBot.Send(message.Sender, err.Error())
	}
	valueRange := in.ToValueRange()
	response, err := client.Spreadsheets.Values.Update(sheetID, sheetRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		TelegramBot.Send(message.Sender, err.Error())
	}
	//TODO: Save input state to redis
	msg := fmt.Sprintf("Updated Row %v, Updated Column %v", response.UpdatedRows, response.UpdatedColumns)
	TelegramBot.Send(message.Sender, msg)
}

func HandleTest(message *tb.Message) {
	if Client() == nil {
		TelegramBot.Send(message.Sender, "Client is nil, authorize first")
	}

	srv, err := sheets.New(Client())
	if err != nil {
		TelegramBot.Send(message.Sender, err.Error())
	}

	sheetID := os.Getenv("SPREADSHEET_ID")
	sheetRange := "ELANQIST0609_1137757232!A6:F"
	resp, err := srv.Spreadsheets.Values.Get(sheetID, sheetRange).Do()
	if err != nil {
		TelegramBot.Send(message.Sender, err.Error())
	}

	if len(resp.Values) == 0 {
		TelegramBot.Send(message.Sender, "No data found.")
	}

	for _, row := range resp.Values {
		msg := fmt.Sprintf("%s %s", row[0], row[1])
		TelegramBot.Send(message.Sender, msg)
	}

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
