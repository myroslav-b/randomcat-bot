package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const cBotHostName = "https://6d9b94b9babf.ngrok.io"

//TSendPhoto contains information for sending photos via API Telegram
type TSendPhoto struct {
	ChatID int    `json:"chat_id"`
	URL    string `json:"photo"`
}

//TSendMessage contains information for sending message via API Telegram
type TSendMessage struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
}

//TChat contains information about Telegram chat
type TChat struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

//TMessage contains information about the message contained in the hook
type TMessage struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	Chat      TChat
}

//TUpdate contains information obtained by unmarshaling of the hook body
type TUpdate struct {
	UpdateID int `json:"update_id"`
	//Message  interface{} `json:"message"`
	Message TMessage
}

func apiSetWebhook(bot TBot) error {
	url := strings.Join([]string{cBotHostName, bot.Token, "getcat"}, "/")
	req, err := http.NewRequest("GET", strings.Join([]string{bot.apiURL, bot.Token, "/setWebhook", "?", "url=", url}, ""), nil)
	if err != nil {
		return errors.Wrap(err, "SetWebhook error (apiSetWebhook)")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "SetWebhook error (apiSetWebhook)")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.Wrap(errors.New(resp.Status), "SetWebhook error (apiSetWebhook)")
	}
	return nil
}

func apiSendMessage(bot TBot, message TSendMessage) error {
	jsonMessage, err := json.Marshal(message)
	//log.Println(string(jsonMessage))

	r := bytes.NewReader(jsonMessage)
	req, err := http.NewRequest("POST", strings.Join([]string{bot.apiURL, bot.Token, "/sendMessage"}, ""), r)
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		return errors.Wrap(err, "SendMessage error: unable to create request (apiSendMessage)")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "SendMessage error: unable to send request (apiSendMessage)")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(strings.Join([]string{"SendMessage error (apiSendMessage):", resp.Status}, ""))
	}
	return nil
}

func apiSendPhoto(bot TBot, photo TSendPhoto) error {
	jsonPhoto, err := json.Marshal(photo)
	//log.Println(string(jsonPhoto))

	r := bytes.NewReader(jsonPhoto)
	req, err := http.NewRequest("POST", strings.Join([]string{bot.apiURL, bot.Token, "/sendPhoto"}, ""), r)
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		return errors.Wrap(err, "SendPhoto error: unable to create request (apiSendPhoto)")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "SendPhoto error: unable to send request (apiSendPhoto)")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(strings.Join([]string{"SendPhoto error (apiSendPhoto): ", resp.Status}, ""))
	}
	return nil
}

func sendStartMessage(bot TBot, update TUpdate) {
	var message TSendMessage
	message.ChatID = update.Message.Chat.ID
	message.Text = "Hello! Y'm Randomcat_bot"
	err := apiSendMessage(bot, message)
	if err != nil {
		log.Println(errors.Wrap(err, "Error send a start message (sendStartMessage)"))
	} else {
		log.Println("Bot sent a Start message")
	}
}

func sendHelpMessage(bot TBot, update TUpdate) {
	var message TSendMessage
	message.ChatID = update.Message.Chat.ID
	message.Text = "/start - start message\n/help - this help\ncat or /cat- send Cat photo"
	err := apiSendMessage(bot, message)
	if err != nil {
		log.Println(errors.Wrap(err, "Error send a help message (sendHelpNessage)"))
	} else {
		log.Println("Bot sent a Help message")
	}
}

func sendCatMessage(bot TBot, update TUpdate) {
	var message TSendMessage
	message.ChatID = update.Message.Chat.ID
	message.Text = "Send Cat"
	log.Println(<-bot.ch)
	err := apiSendMessage(bot, message)
	if err != nil {
		log.Println(errors.Wrap(err, "Error send a cat message (sendCatMessage)"))
	} else {
		log.Println("Bot send a Cat message")
	}
}

func sendCatPhoto(bot TBot, update TUpdate) {
	var photo TSendPhoto
	photo.ChatID = update.Message.Chat.ID
	photo.URL = <-bot.ch
	//log.Println(<-bot.ch)
	err := apiSendPhoto(bot, photo)
	if err != nil {
		log.Println(errors.Wrap(err, "Error send a cat photo (sendCatPhoto)"))
	} else {
		log.Println("Bot sent a Cat photo")
	}
}

func sendInvalidMessage(bot TBot, update TUpdate) {
	var message TSendMessage
	message.ChatID = update.Message.Chat.ID
	message.Text = "Invalid Message"
	err := apiSendMessage(bot, message)
	if err != nil {
		log.Println(errors.Wrap(err, "Error send a invalid message (sendInvalidMessage)"))
	} else {
		log.Println("bot sent a Invalid message")
	}
}

func (bot TBot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body []byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.Wrap(err, "Error reading request body (ServeHTTP)"))
		return
	}
	//log.Printf("%s", body)
	var update TUpdate
	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Println(errors.Wrap(err, "Error unmarshaling request body (ServeHTTP)"))
		return
	}
	//log.Println(update)

	switch update.Message.Text {
	case "/start":
		sendStartMessage(bot, update)
	case "/help":
		sendHelpMessage(bot, update)
	case "cat", "/cat":
		sendCatPhoto(bot, update)
	default:
		sendInvalidMessage(bot, update)
	}

}

func (bot TBot) run() {
	patchGetCat := strings.Join([]string{"", bot.Token, "getcat"}, "/")
	hostAndPort := strings.Join([]string{bot.Address, bot.Port}, ":")
	log.Println(patchGetCat)
	log.Println(hostAndPort)

	defer close(chanCat)
	go catGenerator(chanCat)

	err := apiSetWebhook(bot)
	if err != nil {
		log.Println(err)
		return
	}

	//http.HandleFunc(patchGetCat, handlerGetCat)
	http.Handle(patchGetCat, bot)
	log.Println("Starting bot server")
	http.ListenAndServe(hostAndPort, nil)

}
