package constants

import "time"

const HashedKey = "793ee2d0d52ecabc6017813e9397a659f12f622d390968147d81a842d135560f"

const ExpireTime = 1628619190

// PingPeriod Send pings to peer with this period. Must be less than pongWait.
const PingPeriod = 15 * time.Second

// PongWait Time allowed to read the next pong message from the peer.
const PongWait = 20 * time.Second

// PositionsInitializingCoolDown Time for which fills are requested to prevent contradiction with fill ws channel
const PositionsInitializingCoolDown = 15 * time.Second

// TimeoutToCheckForSubscriptions Timeout to check if channels are subscribed else restart
const TimeoutToCheckForSubscriptions = 15 * time.Second

const OrderIDEncryptionKey = "EF4359D8D580AA4F7F036D6F04FC6A94"

const OrderIDEncryptionTweak = "D8E7920AFA330A73"

const ClientOrderIDPrefix = "DTM"

const ClientOrderIDSuffixLength = 4
