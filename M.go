package main

import (
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/fpe"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println(fpe.GenerateRandomClOrdID())
	orderId := "85724867353245"
	clOrderId := fpe.GenerateClOrdIDFromOrdID(orderId)

	startTime := time.Now().UnixNano()
	verifyClientID(clOrderId)
	fmt.Printf("Time take to verify: %d\n", time.Now().UnixNano()-startTime)

	fmt.Println("Client Id:", clOrderId)
	oNew := fpe.GenerateOrdIDFromClOrdID(clOrderId)
	fmt.Println("decrypted order id:", oNew)
}

func verifyClientID(clientID string) bool {
	if !strings.HasPrefix(clientID, constants.ClientOrderIDPrefix) {
		return false
	}

	clID := strings.TrimPrefix(clientID, constants.ClientOrderIDPrefix)
	if len(clientID) < 2+constants.ClientOrderIDSuffixLength {
		return false
	}

	clID = clID[:len(clID)-constants.ClientOrderIDSuffixLength]
	if _, err := strconv.Atoi(clID); err != nil {
		return false
	}

	return true
}
