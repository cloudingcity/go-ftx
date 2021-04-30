package stream

import (
	"encoding/json"
	"math"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	sec, nsec := math.Modf(f)
	t.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}

type General struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type Error struct {
	Type string `json:"type"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Pong struct {
	Type string `json:"op"`
}

type OrderBook struct {
	General
	Data struct {
		Bids     [][]float64 `json:"bids"`
		Asks     [][]float64 `json:"asks"`
		Time     *Time       `json:"time"`
		Checksum int         `json:"checksum"`
		Action   string      `json:"action"`
	} `json:"data"`
}

type Trade struct {
	General
	Data []struct {
		ID          int       `json:"id"`
		Liquidation bool      `json:"liquidation"`
		Price       float64   `json:"price"`
		Side        string    `json:"side"`
		Size        float64   `json:"size"`
		Time        time.Time `json:"time"`
	} `json:"data"`
}

type Ticker struct {
	General
	Data struct {
		Bid     float64 `json:"bid"`
		Ask     float64 `json:"ask"`
		BidSize float64 `json:"bidSize"`
		AskSize float64 `json:"askSize"`
		Last    float64 `json:"last"`
		Time    *Time   `json:"time"`
	} `json:"data"`
}

type Fills struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Data    struct {
		Fee       float64   `json:"fee"`
		FeeRate   float64   `json:"feeRate"`
		Future    string    `json:"future"`
		ID        int       `json:"id"`
		Liquidity string    `json:"liquidity"`
		Market    string    `json:"market"`
		OrderID   int       `json:"orderId"`
		TradeID   int       `json:"tradeId"`
		Price     float64   `json:"price"`
		Side      string    `json:"side"`
		Size      float64   `json:"size"`
		Time      time.Time `json:"time"`
		Type      string    `json:"type"`
	} `json:"data"`
}

type Orders struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Data    struct {
		ID            int     `json:"id"`
		ClientID      string  `json:"clientId"`
		Market        string  `json:"market"`
		Type          string  `json:"type"`
		Side          string  `json:"side"`
		Size          float64 `json:"size"`
		Price         float64 `json:"price"`
		ReduceOnly    bool    `json:"reduceOnly"`
		Ioc           bool    `json:"ioc"`
		PostOnly      bool    `json:"postOnly"`
		Status        string  `json:"status"`
		FilledSize    float64 `json:"filledSize"`
		RemainingSize float64 `json:"remainingSize"`
		AvgFillPrice  float64 `json:"avgFillPrice"`
	} `json:"data"`
}
