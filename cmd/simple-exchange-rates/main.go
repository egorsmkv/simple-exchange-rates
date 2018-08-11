package main

import (
	"github.com/integrii/flaggy"
	"github.com/tgbotslab/simple-exchange-rates"
)

const version = "0.1"

var token string

// 697543130:AAH4cXtzwmNmilPqD217BKWHs8ZImMvthps

func init() {
	flaggy.SetName("Simple Exchange Rates")
	flaggy.SetDescription("This bot shows the rate of one currency to another currency.")

	flaggy.String(&token, "t", "token", "telegram bot token")

	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {
	bot.Init(token)
}
