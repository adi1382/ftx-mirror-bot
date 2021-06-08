package client

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

func NewSubClient(
	apiKey, apiSecret string,
	leverageUpdateDuration, balanceUpdateDuration, calibrationDuration int64,
	isCopyLeverage, isBalanceProportional bool,
	fixedProportion float64,
	subRoutineCloser chan int, wg *sync.WaitGroup,
	host *Host) *Sub {

	c := Sub{
		client: newClient(apiKey, apiSecret, leverageUpdateDuration, balanceUpdateDuration, subRoutineCloser, wg), hostClient: host,
	}
	c.calibrationDuration = calibrationDuration
	c.isCopyLeverage = isCopyLeverage
	c.isBalanceProportional = isBalanceProportional
	c.fixedProportion = fixedProportion
	c.hostMessageUpdates = make(chan []byte, 100)
	return &c
}

type Sub struct {
	client                *client
	hostClient            *Host
	hostMessageUpdates    chan []byte
	subRoutineCloser      chan int
	wg                    *sync.WaitGroup
	calibrationDuration   int64
	isCopyLeverage        bool
	isBalanceProportional bool
	fixedProportion       float64
}

func (s *Sub) Initialize() {

	//TODO: FILLS ADJUSTMENT
	s.setSymbolInformationFromHost()
	s.client.initialize()
	s.hostClient.SubscribeToHostUpdates(s.hostMessageUpdates)
}

func (s *Sub) setSymbolInformationFromHost() {
	s.client.symbolInfoLock.Lock()
	defer s.client.symbolInfoLock.Unlock()

	symbolInformation := s.hostClient.SymbolInformation()

	s.client.symbolsInfo = make(map[string]symbolInfo, 1000)

	for k := range symbolInformation {
		s.client.symbolsInfo[k] = symbolInformation[k]
	}
}

func (s *Sub) adjustedSize(size float64, symbol string) float64 {
	var adjustedSize float64
	if s.isBalanceProportional {
		adjustedSize = size * (s.client.totalCollateral.Load() / s.hostClient.TotalCollateral())
	} else {
		adjustedSize = size * s.fixedProportion
	}

	return s.roundSize(adjustedSize, symbol)
}

func (s *Sub) roundSize(adjustedSize float64, symbol string) float64 {
	s.client.symbolInfoLock.Lock()
	defer s.client.symbolInfoLock.Unlock()

	rounded := math.Round(adjustedSize/s.client.symbolsInfo[symbol].SizeIncrement) * s.client.symbolsInfo[symbol].SizeIncrement
	formatted, err := strconv.ParseFloat(fmt.Sprintf("%.10f", rounded), 64)
	if err != nil {
		return rounded
	}
	return formatted
}
