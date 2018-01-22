package ticker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

// GLOBAL VARS
var tickers = make(map[string]Ticker)
var readsAll = make(chan *readAllOp)
var readsOne = make(chan *readOneOp)
var writes = make(chan *writeOp, 20)
var pairs = map[string]string{
	"BTCUSD": "USDT-BTC",
	"ETHUSD": "USDT-ETH",
	"ETHBTC": "BTC-ETH",
	"LTCUSD": "USDT-LTC",
	"LTCBTC": "BTC-LTC",
	"XRPUSD": "USDT-XRP",
	"XRPBTC": "BTC-XRP",
	"ZECUSD": "USDT-ZEC",
	"ZECBTC": "BTC-ZEC",
	"XMRUSD": "USDT-XMR",
	"XMRBTC": "BTC-XMR",
	"DASHUSD": "USDT-DASH",
	"DASHBTC": "BTC-DASH",
	"BCHUSD": "USDT-BCC",
	"BCHBTC": "BTC-BCC",
	"ETCUSD": "USDT-ETC",
	"ETCBTC": "BTC-ETC",
}

// Get a ticker from Kraken api
func getTicker(pair string) (aTickerResponse TickerResponse, err error) {

	httpCLI := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	url := fmt.Sprintf("https://bittrex.com/api/v1.1/public/getticker?market=%s", pair)

	// try to get kraken ticker
	resp, err := httpCLI.Get(url)
	if err != nil {
		return TickerResponse{}, err
	}

	// make sure the body of the response is closed after func returns
	defer resp.Body.Close()

	// try to read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TickerResponse{}, err
	}

	// Unmarshal the json
	tickerResponse := TickerResponse{}
	err = json.Unmarshal(body, &tickerResponse)
	if err != nil {
		return TickerResponse{}, err
	}

	return tickerResponse, nil
}


// STATE
func InitState() {
	var state = make(map[string]Ticker)

	for {
		select {
		case read := <-readsAll:
			read.resp <- state

		case read := <-readsOne:
			read.resp <- state[read.key]

		case write := <-writes:
			//if write.key == "BTCUSD" {
			//
			//	fmt.Printf("%s %s\n", write.val.LastPrice, time.Unix(write.val.Timestamp, 0))
			//}
			state[write.key] = write.val
			write.resp <- true
		}
	}
}

// WRITE new tickers
func write(pair string, result TickerResult) {
	ticker := Ticker{
		LastPrice: strconv.FormatFloat(result.Last, 'f', 8, 64),
		BestBid:   strconv.FormatFloat(result.Bid, 'f', 8, 64),
		BestAsk:   strconv.FormatFloat(result.Ask, 'f', 8, 64),
		Timestamp: time.Now().Unix(),
	}

	write := &writeOp{
		key:  pair,
		val:  ticker,
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
				tickerResponse, err := getTicker(xchK)
				if err == nil && tickerResponse.Success {
					write(sg3K, tickerResponse.Result)
				}
			}
		}(sg3Key, xchKey)
	}
}
