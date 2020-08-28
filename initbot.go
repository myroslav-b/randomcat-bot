package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	"gopkg.in/ini.v1"

	"github.com/pkg/errors"
)

const cIniFileName = "config.ini"

func matchToken(token string) bool {
	const cRegexpToken = `^[0-9]{10}:[a-zA-Z0-9_-]{35}$`
	b, err := regexp.MatchString(cRegexpToken, token)
	if !b || err != nil {
		return false
	}
	return true
}

func matchIP(ip string) bool {
	const cRegexpIP = `^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`
	b, err := regexp.MatchString(cRegexpIP, ip)
	if !b || err != nil {
		return false
	}
	return true
}

func matchDN(dn string) bool {
	const cRegexpDN = `^(([A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9]|[A-Za-z0-9]){1,63}\.)+[A-Za-z]{2,6}$`
	b, err := regexp.MatchString(cRegexpDN, dn)
	if !b || err != nil {
		return false
	}
	return true
}

func matchAddress(address string) bool {
	return matchIP(address) || matchDN(address)
}

func matchPort(port string) bool {
	//const cRegexpPort = `^[1-9][0-9]{0,4}$`
	const cRegexpPort = `^6[0-5][0-5][0-3][0-5]$|^[1-5][0-9]{0,4}$|^[1-9][0-9]{0,3}$`
	b, err := regexp.MatchString(cRegexpPort, port)
	if !b || err != nil {
		return false
	}
	return true
}

func iniInit(bot *TBot, config interface{}) bool {
	botSave := bot
	cnf, err := ini.Load(config)
	if err != nil {
		log.Println(errors.Wrap(err, "Error reading configuration file (iniInit)"), " | The configuration file will not be used")
		return false
	}
	err = cnf.MapTo(bot)
	if err != nil {
		log.Println(errors.Wrap(err, "Configuration parsing error (iniInit)"), " | The configuration file will not be used")
		return false
	}
	if bot.LogWriterName == "" {
		log.Println("Log file name error for configuration (iniInit). The log file name from the configuration file will not be used")
		bot.LogWriterName = botSave.LogWriterName
	}

	if !matchToken(bot.Token) && bot.Token != "" {
		log.Println("Tocken error for configuration (iniInit). The token from the configuration file will not be used")
		bot.Token = botSave.Token
	}
	if !matchPort(bot.Port) {
		log.Println("Port error for configuration (iniInit). The port from the configuration file will not be used")
		bot.Port = botSave.Port
	}
	if !matchAddress(bot.Address) && bot.Address != "localhost" {
		log.Println("Address error for configuration (iniInit). The address from the configuration file will not be used")
		bot.Address = botSave.Address
	}
	return true
}

func flagInit(bot *TBot) bool {
	logFileName := bot.LogWriterName
	token := bot.Token
	address := bot.Address
	port := bot.Port

	flag.StringVar(&logFileName, "logfile", bot.LogWriterName, "Name of log file. If logfile is missing or empty - log is written in stdout")
	flag.StringVar(&token, "token", bot.Token, "Token for the bot. Mandatory parameter. Cannot be empty or missing")
	flag.StringVar(&address, "address", bot.Address, "Address (IP or DN) for the bot. If missing: 127.0.0.1")
	flag.StringVar(&port, "port", bot.Port, "Port for the bot. If missing: 443")
	flag.Parse()

	if logFileName != "stdout" {
		file, err := os.Create(logFileName)
		if err != nil {
			log.Println(errors.Wrap(err, "Reading start options: unable to open log file (flagInit)"), " | The parameter is ignored")
		} else {
			bot.LogWriterName = logFileName
			bot.logWriter = file
		}
	}
	//log.SetOutput(bot.logWriter)

	if matchAddress(address) {
		bot.Address = address
	} else {
		log.Println("Reading start options: wrong address (flagInit). The parameter is ignored")
	}

	if matchToken(token) {
		bot.Token = token
	} else {
		log.Println("Reading start options: wrong token (flagInit). The parameter is ignored")
	}

	if matchPort(port) {
		bot.Port = port
	} else {
		log.Println("Reading start options: wrong port (flagInit). The parameter is ignored")
	}

	return true
}

func (bot *TBot) init() bool {
	iniInit(bot, cIniFileName)
	flagInit(bot)
	if bot.Token == "" {
		log.Println("Token is empty. Program stopped")
		return false
	}
	log.Println("Next logs will by written to ", bot.LogWriterName)
	log.SetOutput(bot.logWriter)
	return true
}
