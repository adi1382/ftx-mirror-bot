package tools

import (
	"bufio"
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/constants"
	"os"
	"strconv"
	"strings"
	"time"
)

func RoundFloat(f float64) float64 {
	r, _ := strconv.ParseFloat(fmt.Sprintf("%.10f", f), 64)
	return r
}

func RoundFloatPointer(f *float64) {
	*f, _ = strconv.ParseFloat(fmt.Sprintf("%.10f", *f), 64)
}

//VerifyClientID verifies if the clientID is placed by mirror bot or not
func VerifyClientID(clientID string) bool {
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

func EnterToExit(errMessage string) {
	fmt.Println(errMessage)
	fmt.Print("\n\nPress 'Enter' to exit")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

func CheckIfLicenseExpired() {
	if time.Now().Unix() > constants.ExpireTime {
		EnterToExit("License Expired!!")
	}
}

func ExitIfLicenseExpired() {
	for {
		if time.Now().Unix() > constants.ExpireTime {
			fmt.Println("License Expired!")
			time.Sleep(time.Second * 10)
			EnterToExit("License Expired!!")
		}
		time.Sleep(5 * time.Second)
	}
}
