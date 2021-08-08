// +build windows

package tools

import (
	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/wmic"
)

func CheckLicense() bool {
	if constants.HashedKey == wmic.GetHashedKey() {
		return true
	}
	return false
}
