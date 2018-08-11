package bot

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/openprovider/ecbrates"
	"github.com/shopspring/decimal"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	howToUse            = "Type an exchange way. For example: 200 usd in uah."
	somethingWentWrong  = "Sorry, but something went wrong. Try again later."
	wrongQuestion       = "Wrong question to the bot."
	caanotGetExRates    = "Cannot to get exchange rates."
	cannotParseAmount   = "Sorry, but the bot cannot parse the amount of your question"
	cannotParseCurrency = "Sorry, but the bot cannot parse a currency from your question"
	cannotConvert       = "Sorry, but the bot cannot make this exchange"
	currenciesSame      = "Currencies are the same"
)

type parsedQuestion struct {
	Amount float64
	From   ecbrates.Currency
	To     ecbrates.Currency
}

func Init(token string) {
	log.Println("Simple Exchange Rates is started")

	// Connect to Telegram bot API
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// It is a shortcut to send a reply
	reply := func(b *tb.Bot, m *tb.Message, msg string) {
		b.Send(m.Sender, msg)
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		err, result := parseQuestion(m.Text)
		if err != nil {
			reply(b, m, err.Error())
			return
		}

		reply(b, m, result)
	})

	b.Handle("/start", func(m *tb.Message) {
		reply(b, m, howToUse)
	})

	b.Start()
}

func parseQuestion(text string) (error, string) {
	r, err := ecbrates.New()
	if err != nil {
		return errors.New(caanotGetExRates), ""
	}

	upperText := strings.ToUpper(text)
	parts := strings.Split(upperText, " ")

	if len(parts) != 4 {
		return errors.New(wrongQuestion), ""
	}

	val, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return errors.New(cannotParseAmount), ""
	}

	var res parsedQuestion

	res.Amount = val
	res.From = ecbrates.Currency(parts[1])
	res.To = ecbrates.Currency(parts[3])

	if res.Amount < 0 {
		return errors.New(cannotParseCurrency), ""
	}

	if !res.From.IsValid() || !res.To.IsValid() {
		return errors.New(cannotParseCurrency), ""
	}

	if res.From == res.To {
		return errors.New(currenciesSame), ""
	}

	value, err := r.Convert(res.Amount, res.From, res.To)
	if err != nil {
		return errors.New(cannotConvert), ""
	}

	prettyVal := decimal.NewFromFloat(value).Round(3).String()

	return nil, prettyVal
}
