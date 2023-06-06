package fas

type OrderState int

const (
	OPEN OrderState = iota
	CANCELED
	FILLED
)

const LONG = true
const SHORT = false

type TransferType int

const (
	WITHDRAW TransferType = iota
	DEPOSIT
	TRANSFER
)
