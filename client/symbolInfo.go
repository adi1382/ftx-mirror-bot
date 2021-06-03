package client

type symbolInfo struct {
	symbol         string
	enabled        bool
	priceIncrement float64
	sizeIncrement  float64
	marketType     string //spot, perpetual or futures
}

func (c *Client) updateSymbolInfo() {
	marketResponse := c.getAllMarkets()
	futuresResponse := c.getAllFuturesList()

	c.symbolsInfo = make([]*symbolInfo, 0, 600)

	for i := range *futuresResponse {
		c.symbolsInfo = append(c.symbolsInfo, &symbolInfo{
			symbol:         (*futuresResponse)[i].Name,
			enabled:        (*futuresResponse)[i].Enabled,
			priceIncrement: (*futuresResponse)[i].PriceIncrement,
			sizeIncrement:  (*futuresResponse)[i].SizeIncrement,
			marketType:     (*futuresResponse)[i].Type,
		})
	}

	for i := range *marketResponse {
		if (*marketResponse)[i].Type == "spot" {
			c.symbolsInfo = append(c.symbolsInfo, &symbolInfo{
				symbol:         (*marketResponse)[i].Name,
				enabled:        (*marketResponse)[i].Enabled,
				priceIncrement: (*marketResponse)[i].PriceIncrement,
				sizeIncrement:  (*marketResponse)[i].SizeIncrement,
				marketType:     "spot",
			})
		}
	}

}
