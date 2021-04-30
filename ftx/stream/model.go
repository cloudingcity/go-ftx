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
