package configuration

type Config struct {
	Settings    Settings              `json:"Settings"`
	HostAccount HostAccount           `json:"HostAccount"`
	SubAccounts map[string]SubAccount `json:"SubAccounts"`
}

type Settings struct {
	RatioUpdateRate int `json:"RatioUpdateRate"`
	CalibrationRate int `json:"CalibrationRate"`
}

type HostAccount struct {
	ApiKey string `json:"ApiKey"`
	Secret string `json:"secret"`
}

type SubAccount struct {
	Enabled           bool   `json:"Enabled"`
	BalanceProportion bool   `json:"BalanceProportion"`
	FixedProportion   int    `json:"FixedProportion"`
	CopyLeverage      bool   `json:"CopyLeverage"`
	ApiKey            string `json:"ApiKey"`
	Secret            string `json:"secret"`
}
