package constants

import "time"

// PingPeriod Send pings to peer with this period. Must be less than pongWait.
const PingPeriod = 15 * time.Second

// PongWait Time allowed to read the next pong message from the peer.
const PongWait = 20 * time.Second

// PositionsInitializingCoolDown Time for which fills are requested to prevent contradiction with fill ws channel
const PositionsInitializingCoolDown = 15 * time.Second
