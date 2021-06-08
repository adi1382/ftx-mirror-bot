package client

import (
	"math"

	"github.com/adi1382/ftx-mirror-bot/tools"
)

func (c *client) setFillsAdjuster() {
	c.symbolInfoLock.Lock()
	c.fillsAdjusterLock.Lock()
	defer c.symbolInfoLock.Unlock()
	defer c.fillsAdjusterLock.Unlock()

	c.fillsAdjuster = make(map[string]float64, 600)

	for i := range c.symbolsInfo {
		c.fillsAdjuster[c.symbolsInfo[i].Symbol] = 0
	}
}

func (c *client) addToFillsAdjuster(market string, fillSize float64) {
	c.fillsAdjusterLock.Lock()
	defer c.fillsAdjusterLock.Unlock()

	c.fillsAdjuster[market] = tools.RoundFloat(c.fillsAdjuster[market] + fillSize)
}

func (c *client) removeFromFillsAdjuster(market, side string, fillSize *float64) {
	c.fillsAdjusterLock.Lock()
	defer c.fillsAdjusterLock.Unlock()

	var f float64

	if side == "buy" {
		f = *fillSize
	} else {
		f = -1 * (*fillSize)
	}

	if c.fillsAdjuster[market]*f > 0 {
		if math.Abs(c.fillsAdjuster[market]) > math.Abs(f) {
			c.fillsAdjuster[market] = tools.RoundFloat(c.fillsAdjuster[market] - f)
			f = 0
		} else {
			f -= c.fillsAdjuster[market]
			tools.RoundFloatPointer(&f)
			c.fillsAdjuster[market] = 0
		}
	}

	*fillSize = math.Abs(f)
	tools.RoundFloatPointer(fillSize)
}
