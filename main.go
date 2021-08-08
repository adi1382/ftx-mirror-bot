package main

import (
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/configuration"
	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/mirror"
	"github.com/adi1382/ftx-mirror-bot/tools"
	"sync"
	"time"

	"github.com/adi1382/ftx-mirror-bot/client"
)

//TODO: check all combinations of orders
//TODO: reduce only orders -> check open position size first.
//TODO: market orders in IOC
var (
	subRoutineCloser chan int
	wg               sync.WaitGroup
)

func init() {
	subRoutineCloser = make(chan int, 100)

	if !tools.CheckLicense() {
		fmt.Println("License Invalid")
		time.Sleep(time.Second * 10)
		panic("License invalid")
	}

	fmt.Println(time.Now().Add(30 * 24 * time.Hour).Unix())

	fmt.Printf("Time left to expiration: %d days.\n", ((constants.ExpireTime-time.Now().Unix())/3600)/24)

	tools.CheckIfLicenseExpired()
	go tools.ExitIfLicenseExpired()
}

func main() {

	for {
		config := configuration.ReadConfig()

		hostClient := client.NewHostClient(config.HostAccount.ApiKey,
			config.HostAccount.Secret,
			config.HostAccount.IsFTXSubAccount,
			config.HostAccount.FTXSubAccountName,
			config.Settings.LeverageUpdateDuration,
			config.Settings.CollateralUpdateDuration,
			subRoutineCloser, &wg)

		newMirror := mirror.NewMirrorInstance(subRoutineCloser, &wg)

		newMirror.SetHostClient(hostClient)

		for i := range config.SubAccounts {
			if config.SubAccounts[i].Enabled {
				newMirror.AddSubClient(client.NewSubClient(
					config.SubAccounts[i].ApiKey,
					config.SubAccounts[i].Secret,
					config.SubAccounts[i].IsFTXSubAccount,
					config.SubAccounts[i].FTXSubAccountName,
					config.Settings.LeverageUpdateDuration,
					config.Settings.CollateralUpdateDuration,
					config.Settings.CalibrationRate,
					config.SubAccounts[i].CopyLeverage,
					config.SubAccounts[i].BalanceProportion,
					config.SubAccounts[i].FixedProportion,
					subRoutineCloser, &wg,
					hostClient))
			}
		}

		newMirror.Initialize()
		newMirror.StartMirroring()

		wg.Wait()
		fmt.Println("wait group completed")

		time.Sleep(time.Second * 5)
		fmt.Println("Restarting")
	}

	//n := 0

	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Positions: %v\n", hostClient.ActivePositions())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Orders: %v\n", hostClient.ActiveOrders())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()

	//fmt.Println("$$$$$$$444")
}

//func main() {
//
//	//bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
//	//if _, err := rand.Read(bytes); err != nil {
//	//	panic(err.Error())
//	//}
//
//	//key := hex.EncodeToString(bytes)[0:32] //encode key in bytes to string and keep as secret, put in a vault
//	//fmt.Printf("key to encrypt/decrypt : %s\n", key)
//
//	bytes := []byte("12345678")
//
//	startTime := time.Now().UnixNano()
//	encrypted := encrypt(bytes,"Hello Encrypt")
//	fmt.Printf("Time taken to encrypt: %d\n", time.Now().UnixNano()-startTime)
//
//	fmt.Printf("encrypted : %s\n", encrypted)
//
//	startTime = time.Now().UnixNano()
//	decrypted := decrypt(bytes, encrypted)
//	fmt.Printf("Time taken to decrypt: %d\n", time.Now().UnixNano()-startTime)
//	fmt.Printf("decrypted : %s\n", decrypted)
//}
//
//func encrypt(key []byte, text string) string {
//	// key := []byte(keyText)
//	plaintext := []byte(text)
//
//	block, err := des.NewCipher(key)
//	if err != nil {
//		panic(err)
//	}
//
//	// The IV needs to be unique, but not secure. Therefore it's common to
//	// include it at the beginning of the ciphertext.
//	ciphertext := make([]byte, des.BlockSize+len(plaintext))
//	iv := ciphertext[:des.BlockSize]
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//		panic(err)
//	}
//
//	stream := cipher.NewCFBEncrypter(block, iv)
//	stream.XORKeyStream(ciphertext[des.BlockSize:], plaintext)
//
//	// convert to base64
//	return base64.URLEncoding.EncodeToString(ciphertext)
//}
//
//// decrypt from base64 to decrypted string
//func decrypt(key []byte, cryptoText string) string {
//	ciphertext, _ :=         base64.URLEncoding.DecodeString(cryptoText)
//
//	block, err := des.NewCipher(key)
//	if err != nil {
//		panic(err)
//	}
//
//	// The IV needs to be unique, but not secure. Therefore it's common to
//	// include it at the beginning of the ciphertext.
//	if len(ciphertext) < des.BlockSize {
//		panic("ciphertext too short")
//	}
//	iv := ciphertext[:des.BlockSize]
//	ciphertext = ciphertext[des.BlockSize:]
//
//	stream := cipher.NewCFBDecrypter(block, iv)
//
//	// XORKeyStream can work in-place if the two arguments are the same.
//	stream.XORKeyStream(ciphertext, ciphertext)
//
//	return fmt.Sprintf("%s", ciphertext)
//
//}
