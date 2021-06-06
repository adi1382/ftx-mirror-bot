package client

type symbolInfo struct {
	symbol         string
	enabled        bool
	priceIncrement float64
	sizeIncrement  float64
	marketType     string //spot, perpetual or futures
}

func (c *Client) updateSymbolInfo() {
	c.symbolsInfo = make(map[string]symbolInfo, 600)
	c.updateSymbolInfoForFutures()
	c.updateSymbolInfoForSpot()
}

func (c *Client) updateSymbolInfoForFutures() {
	futuresResponse := c.getAllFuturesList()

	if futuresResponse == nil {
		return
	}

	for i := range *futuresResponse {
		c.symbolsInfo[(*futuresResponse)[i].Name] = symbolInfo{
			symbol:         (*futuresResponse)[i].Name,
			enabled:        (*futuresResponse)[i].Enabled,
			priceIncrement: (*futuresResponse)[i].PriceIncrement,
			sizeIncrement:  (*futuresResponse)[i].SizeIncrement,
			marketType:     (*futuresResponse)[i].Type,
		}
	}
}

func (c *Client) updateSymbolInfoForSpot() {
	marketResponse := c.getAllMarkets()

	if marketResponse == nil {
		return
	}

	for i := range *marketResponse {
		if (*marketResponse)[i].Type == "spot" {
			c.symbolsInfo[(*marketResponse)[i].Name] = symbolInfo{
				symbol:         (*marketResponse)[i].Name,
				enabled:        (*marketResponse)[i].Enabled,
				priceIncrement: (*marketResponse)[i].PriceIncrement,
				sizeIncrement:  (*marketResponse)[i].SizeIncrement,
				marketType:     "spot",
			}
		}
	}
}
