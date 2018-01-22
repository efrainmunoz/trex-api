package orderbook

// MODELS
type Order struct {
	Price    string `json:"price"`
	Volume   string `json:"volume"`
}

type Orderbook struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
	Timestamp int64 `json:"timestamp"`
}

type readAllOp struct {
	resp chan map[string]Orderbook
}

type readOneOp struct {
	key  string
	resp chan Orderbook
}

type writeOp struct {
	key  string
	val  Orderbook
	resp chan bool
}

// Models to deal with the third-party api responses
type order struct {
	Quantity float64 `json:"Quantity"`
	Rate     float64 `json:"Rate"`
}

type OrderbookResult struct {
	Buy  []order `json:"Buy"`
	Sell []order `json:"sell"`
}

type OrderbookResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  OrderbookResult `json:"result"`
}