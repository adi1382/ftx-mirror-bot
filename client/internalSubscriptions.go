package client

func (h *Host) SubscribeToHostNewOrderUpdates(ch chan *order) {
	h.client.subscriptionsToNewOrderUpdatesLock.Lock()
	defer h.client.subscriptionsToNewOrderUpdatesLock.Unlock()

	h.client.subscriptionsToNewOrderUpdates = append(h.client.subscriptionsToNewOrderUpdates, ch)
}

func (c *client) sendNewOrderUpdateToSubscriptions(newOrder order) {
	c.subscriptionsToNewOrderUpdatesLock.Lock()
	defer c.subscriptionsToNewOrderUpdatesLock.Unlock()

	if len(c.subscriptionsToNewOrderUpdates) > 0 {
		for i := range c.subscriptionsToNewOrderUpdates {
			c.subscriptionsToNewOrderUpdates[i] <- &newOrder
		}
	}
}

func (h *Host) SubscribeToHostExistingOrderUpdates(ch chan *order) {
	h.client.subscriptionsToExistingUpdatesLock.Lock()
	defer h.client.subscriptionsToExistingUpdatesLock.Unlock()

	h.client.subscriptionsToExistingOrderUpdates = append(h.client.subscriptionsToExistingOrderUpdates, ch)
}

func (c *client) sendExistingOrderUpdateToSubscriptions(canceledOrder order) {
	c.subscriptionsToExistingUpdatesLock.Lock()
	defer c.subscriptionsToExistingUpdatesLock.Unlock()

	if len(c.subscriptionsToExistingOrderUpdates) > 0 {
		for i := range c.subscriptionsToExistingOrderUpdates {
			c.subscriptionsToExistingOrderUpdates[i] <- &canceledOrder
		}
	}
}
