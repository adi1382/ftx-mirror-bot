package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type Config struct {
	mux sync.RWMutex

	Key    string
	Secret string

	// SubAccountID use Account as needed when rewrite ID
	SubAccountID int
	subAccounts  map[int]SubAccount
}

type SubAccount struct {
	UUID     int
	Nickname string
}

func (p *Config) UseSubAccountID(uuid int) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.SubAccountID = uuid
}

func (p *Config) SubAccount() SubAccount {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.subAccounts[p.SubAccountID]
}

func New(key, secret string, subAccounts ...SubAccount) *Config {
	config := &Config{
		Key:          key,
		Secret:       secret,
		SubAccountID: 0,
	}

	if 0 < len(subAccounts) {
		accounts := make(map[int]SubAccount)
		for i := range subAccounts {
			accounts[subAccounts[i].UUID] = subAccounts[i]
		}
		config.subAccounts = accounts
	}

	return config
}

func (p *Config) Signature(body string) string {
	mac := hmac.New(sha256.New, []byte(p.Secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
