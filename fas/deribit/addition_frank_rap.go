package deribit

type GetInstrumentRequest struct {
	InstrumentName string `json:"instrument_name"`
}

type GetInstrumentResponse struct {
	TickSize                 float64 `json:"tick_size"`
	TakerCommission          float64 `json:"taker_commission"`
	SettlementCurrency       string  `json:"settlement_currency"`
	QuoteCurrency            string  `json:"quote_currency"`
	PriceIndex               string  `json:"price_index"`
	MinTradeAmount           float64 `json:"min_trade_amount"`
	MaxLiquidationCommission float64 `json:"max_liquidation_commission"`
	InstrumentName           string  `json:"instrument_name"`
	FutureType               string  `json:"future_type"`
	CounterCurrency          string  `json:"counter_currency"`
	ContractSize             float64 `json:"contract_size"`
	BaseCurrency             string  `json:"base_currency"`
}
