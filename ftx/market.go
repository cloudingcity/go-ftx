package ftx

import (
	"fmt"
	"net/http"
	"time"
)

type MarketService service

const (
	pathMarkets          = "%s/markets"
	pathMarket           = "%s/markets/%s"
	pathMarketsOrderBook = "%s/markets/%s/orderbook"
	pathMarketsTrades    = "%s/markets/%s/trades"
	pathMarketsCandles   = "%s/markets/%s/candles"
)

// All FTX API docs: https://docs.ftx.com/#get-markets
func (s *MarketService) All() ([]Market, error) {
	u := fmt.Sprintf(pathMarkets, s.client.baseURL)

	var out []Market
	if err := s.client.DoPublic(u, http.MethodGet, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type Market struct {
	Name           string  `json:"name"`
	BaseCurrency   string  `json:"baseCurrency"`
	QuoteCurrency  string  `json:"quoteCurrency"`
	Type           string  `json:"type"`
	Underlying     string  `json:"underlying"`
	Enabled        bool    `json:"enabled"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	Last           float64 `json:"last"`
	PostOnly       bool    `json:"postOnly"`
	PriceIncrement float64 `json:"priceIncrement"`
	SizeIncrement  float64 `json:"sizeIncrement"`
	Restricted     bool    `json:"restricted"`
}

// Get FTX API docs: https://docs.ftx.com/#get-single-market
func (s *MarketService) Get(name string) (*Market, error) {
	u := fmt.Sprintf(pathMarket, s.client.baseURL, name)

	var out Market
	if err := s.client.DoPublic(u, http.MethodGet, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type OrderBook struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

type GetOrderBookOptions struct {
	Depth int `url:"depth"`
}

// GetOrderBook FTX API docs: https://docs.ftx.com/#get-orderbook
func (s *MarketService) GetOrderBook(name string, opts *GetOrderBookOptions) (*OrderBook, error) {
	u := fmt.Sprintf(pathMarketsOrderBook, s.client.baseURL, name)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	var out OrderBook
	if err := s.client.DoPublic(u, http.MethodGet, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type Trade struct {
	Id          int       `json:"id"`
	Liquidation bool      `json:"liquidation"`
	Price       float64   `json:"price"`
	Side        string    `json:"side"`
	Size        float64   `json:"size"`
	Time        time.Time `json:"time"`
}

type GetTradesOptions struct {
	Limit     int   `url:"limit"`
	StartTime int64 `url:"start_time"`
	EndTime   int64 `url:"end_time"`
}

// GetTrades FTX API docs: https://docs.ftx.com/#get-trades
func (s *MarketService) GetTrades(name string, opts *GetTradesOptions) ([]Trade, error) {
	u := fmt.Sprintf(pathMarketsTrades, s.client.baseURL, name)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	var out []Trade
	if err := s.client.DoPublic(u, http.MethodGet, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type Candle struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	StartTime time.Time `json:"startTime"`
	Volume    float64   `json:"volume"`
}

type GetHistoricalPrices struct {
	Resolution int   `url:"resolution"`
	Limit      int   `url:"limit"`
	StartTime  int64 `url:"start_time"`
	EndTime    int64 `url:"end_time"`
}

const (
	Resolution15s = 15
	Resolution1m  = 60
	Resolution5m  = 300
	Resolution15m = 900
	Resolution1h  = 3600
	Resolution4h  = 14400
	Resolution1d  = 86400
)

// GetHistoricalPrices FTX API docs: https://docs.ftx.com/#get-historical-prices
func (s *MarketService) GetHistoricalPrices(name string, opts *GetHistoricalPrices) ([]Candle, error) {
	u := fmt.Sprintf(pathMarketsCandles, s.client.baseURL, name)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	var out []Candle
	if err := s.client.DoPublic(u, http.MethodGet, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
