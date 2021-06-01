package rest

import (
	"time"

	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	jsonIter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

const ENDPOINT = "https://ftx.com/api"

var json = jsonIter.ConfigCompatibleWithStandardLibrary

type Client struct {
	Auth *auth.Config

	HTTPC       *fasthttp.Client
	HTTPTimeout time.Duration
}

func New(auth *auth.Config) *Client {
	hc := new(fasthttp.Client)

	return &Client{
		Auth:        auth,
		HTTPC:       hc,
		HTTPTimeout: 5 * time.Second,
	}
}
