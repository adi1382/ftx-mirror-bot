package client

import (
	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/account"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/fills"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"time"
)

type position struct {
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

	accountInformation := new(account.ResponseForInformation)
	fillsResponse := new(fills.Response)

	c.updateAccountInformationAndFills(accountInformation, fillsResponse)

	c.leverage.Store(accountInformation.Leverage)
	c.totalCollateral.Store(accountInformation.Collateral)

	for i := range accountInformation.Positions {
		newPosition := new(position)
		newPosition.Market = accountInformation.Positions[i].Future
		newPosition.Size = accountInformation.Positions[i].NetSize
		newPosition.Side = accountInformation.Positions[i].Side

		c.openPositions = append(c.openPositions, newPosition)
	}
}

func (c *Client) updateAccountInformationAndFills(accountInformation *account.ResponseForInformation, fillsResponse *fills.Response) {
	for {
		accountInformationRestCallTime := time.Now().Unix()
		accountInformation = c.getAccountInformation()
		fillsResponse = c.getFills(constants.PositionsInitializingCoolDown)
		if c.areAnyFillsAfterAccountInformationCall(fillsResponse, accountInformationRestCallTime) {
			continue
		}
		break
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
	}

	c.fillsForPositionInitialization = fillsResponse

	return false
}

///////////////////////// BEGIN --> STREAM ORDER FUNCTIONALITIES /////////////////////////

func (c *Client) handleFillUpdateFromStream(newFill *websocket.FillsData) {
	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	if index := c.checkIfPositionAlreadyExistsForSymbol(newFill); index > -1 {
		c.updateExistingPosition(newFill, index)
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

}

///////////////////////// END --> STREAM ORDER FUNCTIONALITIES /////////////////////////
