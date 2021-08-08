package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/tools"
	"io/ioutil"
	"os"
)

func ReadConfig() Config {
	var jsonFile *os.File
	var err error

	for {
		jsonFile, err = os.Open(ConfigPath)

		if err == nil {
			break
		} else {
			if os.IsNotExist(err) {
				fmt.Println("Config file does not exists")
				continue
			} else {
				fmt.Println("Could not open config file.")
				tools.EnterToExit("Could not open config file.")
			}
		}
	}

	config := Config{}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println("Unable to read the contents of the configuration file.")
		fmt.Println(err)
		tools.EnterToExit("Unable to read the contents of the configuration file.")
	}

	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Invalid Configuration file")
		tools.EnterToExit("Invalid Configuration file")
	}

	return config
}
