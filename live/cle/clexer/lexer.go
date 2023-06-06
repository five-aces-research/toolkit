package clexer

import (
	"strconv"
	"strings"
)

type TokenType int

/*
Grammatic kinda:

[order] ~ [Ticker|Variable] ~ [AMOUNT] ~ ?[MarketPrice] [Price]

ORDER = { SIDE ~ TICKER ~ AMOUNT ~  PRICE}  !
	PRICE ::= {LADERED | FLOAT | PERCENT | DIFF | CLADERED} !
	SIDE ::= {BUY | SELL} !
	AMOUNT ::= {FLOAT | UFLOAT | PERCENT}
	CLADERED ::= { "-l" ~ (high,low) ~ DURATION ~ LADEREDVAR}
	LADERED ::=  { "-l" ~ INTEGER ~ LADEREDVAR }
	LADEREDVAR ::=  {FLOAT ~ FLOAT | UFLOAT ~ UFLOAT | PERCENT ~ PERCENT}


VARIABLE ::= {TEMP | FUNCTION | STRING | AMOUNT}
	STRING ::= { CONSTANT | CONSTANT ~ STRING | LADERED | HIGH | LOW}
	CONSTANT ::= {FLOAT | DIFF | PERCENT | MARKET}
FUNCTION ::={"func" ~ "(" ~ TEMPVAL ~ ")" ~ VARIABLE }
ASSIGN ::= { STRING ~ "="  ~ VARIABLE}

[buy,sell] [ticker] [10%, 0.1, u200] [(market, 38100, d100, 2%), (-(high,low) 3h [d100, 2%]), (-l [38100 38300, d100 d300, 1% 3%])
stop [ticker]


SIDE   VARIABLE   PERCENT FLOAT UFLOAT VARIABLE
[buy,sell] [ticker] [10%, 0.1, u200] [(market, 38100, d100, 2%), (-(high,low) 3h [d100, 2%]), (-l [38100 38300, d100 d300, 1% 3%])
stop [ticker]


cancel side? ticker
cancel side? ticker


funding -position //funding rate der aktuellen positionen
funding -highest 20 //funding der highest 20 coins
stop buy btc-perp [position, u100,0.1] -low 5h

*/

const (
	VARIABLE      TokenType = iota
	TICKER                  // btc-perp
	SIDE                    // buy, sell
	STOP                    // stop
	FLOAT                   // 100 => 100 $ of btc
	UFLOAT                  // 100u  u = unitFloat => buying 100 btc
	PERCENT                 //
	DFLOAT                  // -200 differenceFloat => 200 below/above the price
	ASSIGN                  // = assign to a variable
	FLAG                    // -l -le
	DELETE                  // delete a b c
	FUNC                    // func(a,b,c) creating function
	DURATION                // 4h 1d 30
	LBRACKET                // (
	RBRACKET                // )
	MARKET                  // -market
	SOURCE                  // -high -low -open -close
	CANCEL                  // cancel
	CLOSE                   // close
	FUNDINGPAYS             // fpay | fundingpayments
	POSITION                // -position
	POSITIONORDER           // -po
	FUNDINGRATES            // fundingrates
	READONLY                // -readonly
	REDUCEONLY              // -reduceonly
	SUM                     // -sum
)

type Token struct {
	Type  TokenType
	Value string
}

// Lexer converts an input to a Token Array
func Lexer(input string) (t []Token, err error) {
	in := strings.Split(input, " ") //tokens are seperated with whitespaces. Only exeptions are function, here its forbidden to seperate inside the brakets

	for _, s := range in {
		if len(s) == 0 {
			continue
		}
		last := len(s) - 1
		switch s {
		case "buy", "sell":
			t = append(t, Token{SIDE, s})
		case "stop":
			t = append(t, Token{STOP, s})
		case "delete":
			t = append(t, Token{DELETE, "delete"})
		case "=":
			t = append(t, Token{ASSIGN, ""})
		case "cancel":
			t = append(t, Token{CANCEL, ""})
		case "fpays", "fundingpays":
			t = append(t, Token{FUNDINGPAYS, ""})
		case "frates", "fundingrates":
			t = append(t, Token{FUNDINGRATES, ""})
		case "close":
			t = append(t, Token{CLOSE, ""})
		default:
			switch s[last] {
			case 'h', 'm', 'd':
				v, err := strconv.Atoi(s[:last])
				if err == nil && v > 0 {
					t = append(t, Token{DURATION, s})
					continue
				}
			case '%':
				v, err := strconv.ParseFloat(s[:last], 64)
				if err == nil && v > 0 {
					t = append(t, Token{PERCENT, s[:len(s)-1]})
					continue
				}
			case 'u':
				v, err := strconv.ParseFloat(s[:last], 64)
				if err == nil && v > 0 {
					t = append(t, Token{UFLOAT, s[:len(s)-1]})
					continue
				}
			}

			if len(s) > 6 {
				if s[:5] == "func(" {
					t = append(t, Token{FUNC, "func"}, Token{LBRACKET, "("})
					t = append(t, lexFunc([]byte(s[5:]))...)
					continue
				}
			}

			if s[0] == '-' {
				_, err := strconv.ParseFloat(s[1:], 64)

				if err == nil {
					t = append(t, Token{DFLOAT, s[1:]})
				} else {
					ss := s[1:]
					switch ss {
					case "low", "high", "open", "close":
						t = append(t, Token{SOURCE, ss})
					case "position":
						t = append(t, Token{POSITION, "1.0"})
					case "market":
						t = append(t, Token{MARKET, "1.0"})
					case "reduceonly", "ro":
						t = append(t, Token{REDUCEONLY, ""})
					case "readonly":
						t = append(t, Token{READONLY, ""})
					case "sum":
						t = append(t, Token{SUM, ""})
					case "openorders", "oo":
						t = append(t, Token{POSITIONORDER, ""})
					default:
						t = append(t, Token{FLAG, ss})
					}
				}
				continue
			}

			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				t = append(t, Token{FLOAT, s})
				continue
			}
			t = append(t, lexVariable([]byte(s))...)
		}
	}

	return
}

func lexFunc(s []byte) []Token {
	var temp []byte
	var tk []Token
	for _, v := range s {
		switch v {
		case ')':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			tk = append(tk, Token{VARIABLE, string(temp)})
			temp = []byte("")
		default:
			temp = append(temp, v)
		}
	}

	return tk
}

// LexVariable lexes functions e.g. a(xrp-buy,5,10) and variables
func lexVariable(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{LBRACKET, ""})
			temp = []byte("")
		case ')':
			l, _ := Lexer(string(temp))
			tk = append(tk, l...)
			tk = append(tk, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			temp = append(temp, ' ')
		default:
			temp = append(temp, v)
		}
	}

	if len(temp) != 0 {
		tk = append(tk, Token{VARIABLE, string(temp)})
	}
	return tk
}

func (t Token) Stringer() (out string) {
	switch t.Type {
	case VARIABLE:
		out = "VARIABLE"
	case CANCEL:
		out = "CANCEL"
	case UFLOAT:
		out = "UFLOAT"
	case PERCENT:
		out = "PERCENT"
	case FLAG:
		out = "FLAG"
	case ASSIGN:
		out = "ASSIGN"
	case TICKER:
		return "ticker"
	case SIDE:
		return "side"
	case STOP:
		return "stop"
	case FLOAT:
		return "float"
	case DFLOAT:
		return "DFLOAT"
	case FUNC:
		return "FUNC"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case SOURCE:
		return "SOURCE"
	case CLOSE:
		return "CLOSE"
	case FUNDINGPAYS:
		return "FUNDINGPAYS"
	case POSITION:
		return "POSITION"
	case FUNDINGRATES:
		return "FUNDINGRATES"
	case READONLY:
		return "READONLY"
	case REDUCEONLY:
		return "REDUCEONLY"
	default:
		return t.Value
	}

	return
}
