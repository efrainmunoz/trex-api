package ticker

// MODELS
type Ticker struct {
	LastPrice string `json:"lastprice"`
	BestBid   string `json:"bestbid"`
	BestAsk   string `json:"bestask"`
	Timestamp int64  `json:"timestamp"`
}

type TickerResult struct {
	Bid  float64 `json:"Bid"`
	Ask  float64 `json:"Ask"`
	Last float64 `json:"Last"`
}

type TickerResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Result  TickerResult `json:"result"`
}

type readAllOp struct {
	resp chan map[string]Ticker
}

type readOneOp struct {
	key  string
	resp chan Ticker
}

type writeOp struct {
	key  string
	val  Ticker
	resp chan bool
}

