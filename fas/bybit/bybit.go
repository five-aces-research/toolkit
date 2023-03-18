package bybit

import (
	"log"
	"strings"
)

func categoryTicker(s string) (category, ticker string) {
	s = strings.ToUpper(s)
	ss := strings.Split(s, ".")
	if len(ss) == 1 {
		return "linear", ss[0]
	}
	if len(ss) == 0 {
		log.Panicln(s)
		return
	}
	switch ss[0] {
	case "L":
		category = "linear"
	case "I":
		category = "inverse"
	case "O":
		category = "option"
	case "S":
		category = "spot"
	default:
		log.Println(ss[0])
	}
	return category, ss[1]
}
