package main

import (
	"github.com/efrainmunoz/trex-api/ticker"
	"github.com/efrainmunoz/trex-api/orderbook"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/efrainmunoz/trex-api/trades"
)

// MAIN
func main() {

	// Ticker
	go ticker.InitState()
	go ticker.InitService()

	// Orderbook
	go orderbook.InitState()
	go orderbook.InitService()

	// Trades
	go trades.InitState()
	go trades.InitService()

	// Set api routes
	router := mux.NewRouter()
	router.HandleFunc("/ticker", ticker.GetAll).Methods("GET")
	router.HandleFunc("/ticker/{pair}", ticker.Get).Methods("GET")
	router.HandleFunc("/orderbook", orderbook.GetAll).Methods("GET")
	router.HandleFunc("/orderbook/{pair}", orderbook.Get).Methods("GET")
	router.HandleFunc("/trades", trades.GetAll).Methods("GET")
	router.HandleFunc("/trades/{pair}", trades.Get).Methods("GET")

	// Start the server
	http.ListenAndServe(":8000", router)
}
