package client

type position struct {
	Market string  `json:"market"`
	Size   float64 `json:"size"`
	Side   string  `json:"side"`
}

func (c *Client) initializeAccountInfoAndPositions() {
	//Criteria to call this function
	//1. Websockets should be connected
	//2. There should be no messages pending from Fills channel before calling this function
	//3. There should be no messages pending after calling this functions

	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	c.openPositions = c.openPositions[:0]
	accountInformation := c.getAccountInformation()

	c.leverage.Store(accountInformation.Leverage)
	c.totalCollateral.Store(accountInformation.Collateral)

	for i := range accountInformation.Positions {
		position := new(position)
		position.Market = accountInformation.Positions[i].Future
		position.Size = accountInformation.Positions[i].NetSize
		position.Side = accountInformation.Positions[i].Side

		c.openPositions = append(c.openPositions, position)
	}
}
