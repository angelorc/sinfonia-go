package web3token

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

/*
{
    "signer": "bitsong1zq68dx423frv8yss5skqyg5um97vpefefe2enq",
    "payload": {
        "domain": "test.com",
        "expire_at": 1657293505
    },
    "pub_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "AujBPkj1cLTFNy5OmuP6PCG7ttkLHxCiwUBgG5gNaJZM"
    },
    "signature": "kcQiUQ95qGw/TSJBZ2jKbRqw2ohwnXfNk8XPteT+SAZW3NG3q8198cZx3XdoOca/+9hug1dwd+g0vqZfquPyuQ=="
}

// token
eyJzaWduZXIiOiJiaXRzb25nMXpxNjhkeDQyM2Zydjh5c3M1c2txeWc1dW05N3ZwZWZlZmUyZW5xIiwicGF5bG9hZCI6eyJkb21haW4iOiJ0ZXN0LmNvbSIsImV4cGlyZV9hdCI6MTY1NzI5MzUwNX0sInB1Yl9rZXkiOnsidHlwZSI6InRlbmRlcm1pbnQvUHViS2V5U2VjcDI1NmsxIiwidmFsdWUiOiJBdWpCUGtqMWNMVEZOeTVPbXVQNlBDRzd0dGtMSHhDaXdVQmdHNWdOYUpaTSJ9LCJzaWduYXR1cmUiOiJrY1FpVVE5NXFHdy9UU0pCWjJqS2JScXcyb2h3blhmTms4WFB0ZVQrU0FaVzNORzNxODE5OGNaeDNYZG9PY2EvKzlodWcxZHdkK2cwdnFaZnF1UHl1UT09In0=

// Keplr test
https://jsfiddle.net/dz4kh6nv/1/
*/

const (
	MaxTimeLength = time.Minute * 10
)

type Web3Token struct {
	Signer    string  `json:"signer"`
	Payload   Payload `json:"payload"`
	PubKey    PubKey  `json:"pub_key"`
	Signature string  `json:"signature"`
}

type Payload struct {
	Domain   string `json:"domain"`
	ExpireAt int64  `json:"expire_at"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewWeb3TokenFromBearer(bearer string) (*Web3Token, error) {
	bearerParts := strings.Split(bearer, "Bearer ")

	if len(bearerParts) < 2 {
		return nil, fmt.Errorf("invalid header token")
	}

	token, err := base64.StdEncoding.DecodeString(bearerParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid header token base64")
	}

	var web3token *Web3Token

	err = json.Unmarshal(token, &web3token)
	if err != nil {
		return nil, err
	}

	if ok := web3token.Validate(); !ok {
		return nil, errors.New("wrong signature")
	}

	return web3token, nil
}

func (w3t *Web3Token) SignDoc() []byte {
	payloadBz, err := json.Marshal(w3t.Payload)
	if err != nil {
		return nil
	}

	data := base64.StdEncoding.EncodeToString(payloadBz)

	signDoc := strings.TrimSpace(fmt.Sprintf(`
{
  "chain_id": "",
  "account_number": "0",
  "sequence": "0",
  "fee": {
    "gas": "0",
    "amount": []
  },
  "msgs": [
    {
      "type": "sign/MsgSignData",
      "value": {
        "signer": "%s",
        "data": "%s"
      }
    }
  ],
  "memo": ""
}
`, w3t.Signer, data))

	return []byte(signDoc)
}

func (w3t *Web3Token) GetAddress() sdk.AccAddress {
	appparams.SetAddressPrefixes()
	return w3t.GetPubKey().Address().Bytes()
}

func (w3t *Web3Token) GetPubKey() *secp256k1.PubKey {
	pkBz, err := base64.StdEncoding.DecodeString(w3t.PubKey.Value)
	if err != nil {
		return nil
	}

	return &secp256k1.PubKey{Key: pkBz}
}

func (w3t *Web3Token) GetSignature() []byte {
	sigBz, err := base64.StdEncoding.DecodeString(w3t.Signature)
	if err != nil {
		return nil
	}

	return sigBz
}

func (w3t *Web3Token) GetMsg() []byte {
	var msg map[string]interface{}

	if err := json.Unmarshal(w3t.SignDoc(), &msg); err != nil {
		return nil
	}

	msgBz, _ := json.Marshal(msg)
	return msgBz
}

func (w3t *Web3Token) ValidateSignature() bool {
	return w3t.GetPubKey().VerifySignature(w3t.GetMsg(), w3t.GetSignature())
}

func (w3t *Web3Token) GetDomain() string {
	return w3t.Payload.Domain
}

func (w3t *Web3Token) IsExpired() bool {
	minTime := time.Now()
	maxTime := time.Now().Add(MaxTimeLength)
	parsed := time.Unix(w3t.Payload.ExpireAt, 0)

	return minTime.After(parsed) || maxTime.Before(parsed)
}

func (w3t *Web3Token) Validate() bool {
	return w3t.ValidateSignature() && !w3t.IsExpired()
}
