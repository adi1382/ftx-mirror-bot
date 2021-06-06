package client

import (
	"math"
	"time"

	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/account"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/fills"
	"github.com/adi1382/ftx-mirror-bot/websocket"
)

type Position struct {
	Market string  `json:"market"`
	Size   float64 `json:"size"`
	Side   string  `json:"side"`
}

func (c *Client) initializeAccountInfoAndPositions() {
	// Things to note
	// 1. Positions are made up of fills
	// 2. Any fills are received through WS, older than accountInformation request must be ignored
	// 2. This function fetches fills for last few seconds based on the constant PositionsInitializingCoolDown
	// 3. If any of these fills are received in WS stream, they must be ignored

	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	c.openPositions = c.openPositions[:0]

	accountInformation := c.fetchAccountInformationAndUpdateLastFillTime()

	if accountInformation == nil {
		return
	}

	c.leverage.Store(accountInformation.Leverage)
	c.totalCollateral.Store(accountInformation.Collateral)

	for i := range accountInformation.Positions {
		if accountInformation.Positions[i].NetSize == 0 {
			continue
		}

		newPosition := new(Position)
		newPosition.Market = accountInformation.Positions[i].Future
		newPosition.Size = accountInformation.Positions[i].NetSize
		newPosition.Side = accountInformation.Positions[i].Side

		c.openPositions = append(c.openPositions, newPosition)
	}
}

func (c *Client) fetchAccountInformationAndUpdateLastFillTime() *account.ResponseForInformation {
	for {
		accountInformationRestCallTime := time.Now().Unix()
		accountInformation := c.getAccountInformation()
		if accountInformation == nil {
			return nil
		}
		fillsResponse := c.getFills(constants.PositionsInitializingCoolDown)
		if fillsResponse == nil {
			return accountInformation
		}
		if c.areAnyFillsAfterAccountInformationCall(fillsResponse, accountInformationRestCallTime) {
			time.Sleep(time.Second)
			continue
		}
		return accountInformation
	}
}

func (c *Client) areAnyFillsAfterAccountInformationCall(fillsResponse *fills.Response, accountInformationRestCallTime int64) bool {
	for i := range *fillsResponse {
		if (*fillsResponse)[i].Time.Unix() > accountInformationRestCallTime {
			return true
		}
	}

	if len(*fillsResponse) > 0 {
		c.lastFillUnixTime = (*fillsResponse)[0].Time.Unix()
		c.isPositionCoolDownPeriod.Store(true)
		c.shutDownPositionCoolDownAfter(constants.PositionsInitializingCoolDown)
	}

	return false
}

func (c *Client) shutDownPositionCoolDownAfter(coolDownTime time.Duration) {
	time.AfterFunc(coolDownTime, func() {
		c.isPositionCoolDownPeriod.Store(false)
	})
}

///////////////////////// BEGIN --> STREAM ORDER FUNCTIONALITIES /////////////////////////

func (c *Client) handleFillUpdateFromStream(newFill *websocket.FillsData) {
	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	// c.isPositionCoolDownPeriod is automatically set to off after few seconds after initialization
	if c.isPositionCoolDownPeriod.Load() {
		if newFill.Time.Unix() < c.lastFillUnixTime {
			return
		}
	}

	if val, ok := c.symbolsInfo[newFill.Market]; !ok {
		return
	} else {
		if val.marketType == "spot" {
			return
		}
	}

	if index := c.checkIfPositionAlreadyExistsForSymbol(newFill); index > -1 {
		c.updateExistingPosition(newFill, index)
		c.removePositionIfRequired(index)
		return
	} else {
		index = c.insertNewPosition(newFill)
		c.removePositionIfRequired(index)
		return
	}

}

func (c *Client) checkIfPositionAlreadyExistsForSymbol(newFill *websocket.FillsData) int {
	positionIndex := -1

	for i := range c.openPositions {
		if c.openPositions[i].Market == newFill.Future {
			return i
		}
	}

	return positionIndex
}

func (c *Client) updateExistingPosition(newFill *websocket.FillsData, positionIndex int) {
	var fillSize float64
	if newFill.Side == "buy" {
		fillSize = math.Abs(newFill.Size)
	} else {
		fillSize = -math.Abs(newFill.Size)
	}

	c.openPositions[positionIndex].Size += fillSize

	if c.openPositions[positionIndex].Size >= 0 {
		c.openPositions[positionIndex].Side = "buy"
	} else {
		c.openPositions[positionIndex].Side = "sell"
	}
}

func (c *Client) removePositionIfRequired(positionIndex int) {
	if c.openPositions[positionIndex].Size == 0 {
		c.openPositions = append(c.openPositions[:positionIndex], c.openPositions[positionIndex+1:]...)
	}
}

func (c *Client) insertNewPosition(newFill *websocket.FillsData) int {
	var fillSize float64
	if newFill.Side == "buy" {
		fillSize = math.Abs(newFill.Size)
	} else {
		fillSize = -math.Abs(newFill.Size)
	}

	c.openPositions = append(c.openPositions, &Position{
		Market: newFill.Future,
		Side:   newFill.Side,
		Size:   fillSize,
	})

	return len(c.openPositions) - 1
}

///////////////////////// END --> STREAM ORDER FUNCTIONALITIES /////////////////////////
