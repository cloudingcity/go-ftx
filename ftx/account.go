package ftx

import (
	"fmt"
	"net/http"
)

type AccountService service

const (
	pathAccount         = "%s/account"
	pathPositions       = "%s/positions"
	pathAccountLeverage = "%s/account/leverage"
)

type Account struct {
	BackstopProvider             bool    `json:"backstopProvider"`
	Collateral                   float64 `json:"collateral"`
	FreeCollateral               float64 `json:"freeCollateral"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	Leverage                     float64 `json:"leverage"`
	Liquidating                  bool    `json:"liquidating"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	MakerFee                     float64 `json:"makerFee"`
	MarginFraction               float64 `json:"marginFraction"`
	OpenMarginFraction           float64 `json:"openMarginFraction"`
	TakerFee                     float64 `json:"takerFee"`
	TotalAccountValue            float64 `json:"totalAccountValue"`
	TotalPositionSize            float64 `json:"totalPositionSize"`
	Username                     string  `json:"username"`
	Positions                    []Position
}

// GetInformation FTX API docs: https://docs.ftx.com/#get-account-information
func (s *AccountService) GetInformation() (*Account, error) {
	u := fmt.Sprintf(pathAccount, s.client.baseURL)

	var out Account
	err := s.client.DoPrivate(u, http.MethodGet, nil, &out)
	return &out, err
}

type Position struct {
	Cost                         float64 `json:"cost"`
	EntryPrice                   float64 `json:"entryPrice"`
	EstimatedLiquidationPrice    float64 `json:"estimatedLiquidationPrice,omitempty"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                float64 `json:"unrealizedPnl"`
	CollateralUsed               float64 `json:"collateralUsed,omitempty"`
}

// GetPositions FTX API docs: https://docs.ftx.com/#get-positions
func (s *AccountService) GetPositions() ([]Position, error) {
	u := fmt.Sprintf(pathPositions, s.client.baseURL)

	var out []Position
	err := s.client.DoPrivate(u, http.MethodGet, nil, &out)
	return out, err
}

const (
	Leverage1X   = 1
	Leverage3X   = 3
	Leverage5X   = 5
	Leverage10X  = 10
	Leverage20X  = 20
	Leverage50X  = 50
	Leverage100X = 100
	Leverage101X = 101
)

type RequestLeverage struct {
	Leverage int `json:"leverage"`
}

// SetLeverage FTX API docs: https://docs.ftx.com/#change-account-leverage
func (s *AccountService) SetLeverage(x int) error {
	u := fmt.Sprintf(pathAccountLeverage, s.client.baseURL)

	in := RequestLeverage{Leverage: x}
	return s.client.DoPrivate(u, http.MethodPost, &in, nil)
}
