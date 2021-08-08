package configuration

const ConfigPath = "config/config.json"

type Config struct {
	Settings    Settings              `json:"Settings"`
	HostAccount HostAccount           `json:"HostAccount"`
	SubAccounts map[string]SubAccount `json:"SubAccounts"`
}

type Settings struct {
	LeverageUpdateDuration   int64 `json:"LeverageUpdateDuration"`
	CollateralUpdateDuration int64 `json:"CollateralUpdateDuration"`
	CalibrationRate          int64 `json:"CalibrationRate"`
}

type HostAccount struct {
	ApiKey            string `json:"ApiKey"`
	Secret            string `json:"secret"`
	IsFTXSubAccount   bool   `json:"IsFTXSubAccount"`
	FTXSubAccountName string `json:"FTXSubAccountName"`
}

type SubAccount struct {
	Enabled           bool    `json:"Enabled"`
	BalanceProportion bool    `json:"BalanceProportion"`
	FixedProportion   float64 `json:"FixedProportion"`
	CopyLeverage      bool    `json:"CopyLeverage"`
	ApiKey            string  `json:"ApiKey"`
	Secret            string  `json:"secret"`
	IsFTXSubAccount   bool    `json:"IsFTXSubAccount"`
	FTXSubAccountName string  `json:"FTXSubAccountName"`
}
