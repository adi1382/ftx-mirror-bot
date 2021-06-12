package fpe

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/adi1382/ftx-mirror-bot/constants"
)

var ff3 cipherFF3

func init() {
	key, err := hex.DecodeString(constants.OrderIDEncryptionKey)
	if err != nil {
		panic(err)
	}
	tweak, err := hex.DecodeString(constants.OrderIDEncryptionTweak)
	if err != nil {
		panic(err)
	}

	// Create a new FF1 cipher "object"
	// 10 is the radix/base, and 8 is the tweak length.
	ff3, err = newCipher(10, key, tweak)
	if err != nil {
		panic(err)
	}
}

func GenerateRandomClOrdID() string {
	nBig, err := rand.Int(rand.Reader, big.NewInt(89999999999))
	if err != nil {
		panic(err)
	}
	nBig.Abs(nBig)
	nBig.Add(nBig, big.NewInt(10000000000))
	return GenerateClOrdIDFromOrdID(nBig.Int64())
}

func GenerateClOrdIDFromOrdID(orderID int64) string {
	// Call the encryption function on an example SSN
	clOrderID, err := ff3.Encrypt(strconv.Itoa(int(orderID)))
	if err != nil {
		panic(err)
	}
	clOrderID = constants.ClientOrderIDPrefix + clOrderID + randomString(constants.ClientOrderIDSuffixLength)
	return clOrderID
}

func GenerateOrdIDFromClOrdID(clOrderID string) int64 {
	clOrderID = strings.TrimPrefix(clOrderID, constants.ClientOrderIDPrefix)
	clOrderID = clOrderID[:len(clOrderID)-constants.ClientOrderIDSuffixLength]
	orderID, err := ff3.Decrypt(clOrderID)
	if err != nil {
		panic(err)
	}
	ordID, err := strconv.Atoi(orderID)
	if err != nil {
		panic(err)
	}
	return int64(ordID)
}

func randomString(len int) string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b)[0:len]
}
