package orderbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GLOBAL VARS
var orderbook = make(map[string]Orderbook)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)
var pairs = map[string]string{
	"BTCUSD":  "USDT-BTC",
	"ETHUSD":  "USDT-ETH",
	"ETHBTC":  "BTC-ETH",
	"LTCUSD":  "USDT-LTC",
	"LTCBTC":  "BTC-LTC",
	"XRPUSD":  "USDT-XRP",
	"XRPBTC":  "BTC-XRP",
	"ZECUSD":  "USDT-ZEC",
	"ZECBTC":  "BTC-ZEC",
	"XMRUSD":  "USDT-XMR",
	"XMRBTC":  "BTC-XMR",
	"DASHUSD": "USDT-DASH",
	"DASHBTC": "BTC-DASH",
	"BCHUSD":  "USDT-BCC",
	"BCHBTC":  "BTC-BCC",
	"ETCUSD":  "USDT-ETC",
	"ETCBTC":  "BTC-ETC",
}

// Get an orderbook from Kraken api
func getOrderbook(pair string) (aOrderbookResponse OrderbookResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 2000 * time.Millisecond,
	}

	url := fmt.Sprintf("https://bittrex.com/api/v1.1/public/getorderbook?market=%s&type=both", pair)

	// try to get bittrex orderbook
	resp, err := httpCLI.Get(url)
	if err != nil {
		return OrderbookResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OrderbookResponse{}, err
	}

	// Unmarshal the json
	orderbookResponse := OrderbookResponse{}
	err = json.Unmarshal(body, &orderbookResponse)
	if err != nil {
		return OrderbookResponse{}, err
	}

	return orderbookResponse, nil
}

// STATE
func InitState() {
	var state = make(map[string]Orderbook)

	for {
		select {
		case read := <-readsAll:
			read.resp <- state

		case read := <-readsOne:
			read.resp <- state[read.key]

		case write := <-writes:
			state[write.key] = write.val
			write.resp <- true
		}
	}
}

// WRITE new tickers
func write(pair string, result OrderbookResult) {
	var asks []Order
	var bids []Order

	for _, ask := range result.Sell {
		order := Order{
			Price:  strconv.FormatFloat(ask.Rate, 'f', 8, 64),
			Volume: strconv.FormatFloat(ask.Quantity, 'f', 8, 64)}
		asks = append(asks, order)
	}

	for _, bid := range result.Buy {
		bid := Order{
			Price:  strconv.FormatFloat(bid.Rate, 'f', 8, 64),
			Volume: strconv.FormatFloat(bid.Quantity, 'f', 8, 64)}
		bids = append(bids, bid)
	}

	orderbook := Orderbook{
		Asks:      asks,
		Bids:      bids,
		Timestamp: time.Now().Unix()}

	write := &writeOp{
		key:  pair,
		val:  orderbook,
		resp: make(chan bool)}

	writes <- write
	<-write.resp
}

// Init service
func InitService() {
	for sg3Key, xchKey := range pairs {
		go func(sg3K string, xchK string) {
			ticker := time.NewTicker(time.Millisecond * 1000)
			for range ticker.C {
				orderbookResponse, err := getOrderbook(xchK)
				if err == nil && orderbookResponse.Success {
					write(sg3K, orderbookResponse.Result)
				}
			}
		}(sg3Key, xchKey)
	}
}
