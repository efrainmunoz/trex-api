package trades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

// GLOBAL VARS
var trades = make(map[string]Trades)
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

// Get trades from Kraken api
func getTrades(pair string) (aTradesResponse TradesResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	url := fmt.Sprintf("https://bittrex.com/api/v1.1/public/getmarkethistory?market=%s", pair)

	// try to get kraken trades
	resp, err := httpCLI.Get(url)
	if err != nil {
		return TradesResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TradesResponse{}, err
	}

	// Unmarshal the json
	tradesResponse := TradesResponse{}
	err = json.Unmarshal(body, &tradesResponse)
	if err != nil {
		return TradesResponse{}, err
	}

	return tradesResponse, nil
}


// STATE
func InitState() {
	var state = make(map[string]Trade)

	for {
		select {
		case read := <-readsAll:
			read.resp <- state

		case read := <-readsOne:
			read.resp <- state[read.key]

		case write := <-writes:
			//if write.key == "BTCUSD" {
			//	fmt.Println(write.val)
			//}
			state[write.key] = write.val
			write.resp <- true
		}
	}
}

// WRITE new tickers
func write(sg3Pair string, result TradesResult) {

	var lastTrade trade
	var tradeAction string
	l := len(result)

	if l > 0 {
			lastTrade = result[0]

		if lastTrade.OrderType == "SELL" {
			tradeAction = "sell"
		} else {
			tradeAction = "buy"
		}

		trade := Trade{
			Price: strconv.FormatFloat(lastTrade.Price, 'f', 8, 64),
			Volume: strconv.FormatFloat(lastTrade.Quantity, 'f', 8, 64),
			TradeAction: tradeAction}

		write := &writeOp{
			key:  sg3Pair,
			val: trade,
			resp: make(chan bool)}

		writes <- write
		<-write.resp
	}
}

// Init service
func InitService() {
	for sg3Key, xchKey := range pairs {
		go func(sg3K string, xchK string) {
			ticker := time.NewTicker(time.Millisecond * 1000)
			for range ticker.C {
				tradesResponse, err := getTrades(xchK)
				if err == nil && tradesResponse.Success {
					write(sg3K, tradesResponse.Result)
				}
			}
		}(sg3Key, xchKey)
	}
}
