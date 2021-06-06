package client

type symbolInfo struct {
	Symbol         string
	Enabled        bool
	PriceIncrement float64
	SizeIncrement  float64
	MarketType     string //spot, perpetual or futures
}

func (c *client) updateSymbolInfoForFutures() {
	futuresResponse := c.getAllFuturesList()

	if futuresResponse == nil {
		return
	}

	for i := range *futuresResponse {
		c.symbolsInfo[(*futuresResponse)[i].Name] = symbolInfo{
			Symbol:         (*futuresResponse)[i].Name,
			Enabled:        (*futuresResponse)[i].Enabled,
			PriceIncrement: (*futuresResponse)[i].PriceIncrement,
			SizeIncrement:  (*futuresResponse)[i].SizeIncrement,
			MarketType:     (*futuresResponse)[i].Type,
		}
	}
}

func (c *client) updateSymbolInfoForSpot() {
	marketResponse := c.getAllMarkets()

	if marketResponse == nil {
		return
	}

	for i := range *marketResponse {
		if (*marketResponse)[i].Type == "spot" {
			c.symbolsInfo[(*marketResponse)[i].Name] = symbolInfo{
				Symbol:         (*marketResponse)[i].Name,
				Enabled:        (*marketResponse)[i].Enabled,
				PriceIncrement: (*marketResponse)[i].PriceIncrement,
				SizeIncrement:  (*marketResponse)[i].SizeIncrement,
				MarketType:     "spot",
			}
		}
	}
}
