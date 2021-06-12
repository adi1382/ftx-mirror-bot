package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Config struct {
	Key        string
	secret     string
	SubAccount FTXSubAccount
}

type FTXSubAccount struct {
	IsSubAccount bool
	Name         string
}

func New(key, secret string, isFTXSubAccount bool, FTXSubAccountName string) *Config {
	config := &Config{
		Key:    key,
		secret: secret,
	}

	if isFTXSubAccount {
		if FTXSubAccountName != "" {
			config.SubAccount.IsSubAccount = true
			config.SubAccount.Name = FTXSubAccountName
		}
	}

	return config
}

func (p *Config) Signature(body string) string {
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
