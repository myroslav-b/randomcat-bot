package main

import (
	"io"
	"log"
	"os"
)

//TBot ...
type TBot struct {
	apiURL        string
	LogWriterName string `ini:"logfile"`
	logWriter     io.Writer
	Token         string `ini:"token"`
	Address       string `ini:"address"`
	Port          string `ini:"port"`
	ch            chan string
}

//const cApiTelegramOrgBot = "https://api.telegram.org/bot"
const cCapChanCat = 10

var chanCat chan string = make(chan string, cCapChanCat)

func main() {
	log.SetOutput(os.Stdout)
	defer func() {
		if err := recover(); err != nil {
			log.Println("Global panic: ", err)
		}
	}()

	var bot = TBot{
		"https://api.telegram.org/bot",
		"stdout",
		os.Stdout,
		"",
		"localhost",
		"443",
		chanCat}

	if bot.init() {
		bot.run()
	}
}
