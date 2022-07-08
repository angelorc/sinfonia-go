package web3token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestFromString(t *testing.T) {
	signDoc := strings.TrimSpace(`
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
        "signer": "bitsong1zq68dx423frv8yss5skqyg5um97vpefefe2enq",
        "data": "eyJhZGRyZXNzIjoiYml0c29uZzF6cTY4ZHg0MjNmcnY4eXNzNXNrcXlnNXVtOTd2cGVmZWZlMmVucSIsImRvbWFpbiI6InRlc3QuY29tIiwiZXhwaXJlX2F0IjoxMTExMX0="
      }
    }
  ],
  "memo": ""
}
`)

	pkBz, err := base64.StdEncoding.DecodeString("AujBPkj1cLTFNy5OmuP6PCG7ttkLHxCiwUBgG5gNaJZM")
	require.Nil(t, err)
	require.Len(t, pkBz, 33)

	pub := secp256k1.PubKey{Key: pkBz}

	sigStr, err := base64.StdEncoding.DecodeString("iGf/PXHeDhRw9XL6FW8b7oWrWa/Ed7sVcwqsUo1HwsFkj8HavbpfbHun6WRliK1l+haflxIBzU/dayIY7KwTcA==")
	require.Nil(t, err)

	var msg map[string]interface{}

	err = json.Unmarshal([]byte(signDoc), &msg)
	require.Nil(t, err)

	msgBz, err := json.Marshal(msg)
	require.Nil(t, err)

	fmt.Println(pub.VerifySignature(msgBz, sigStr))
}

func TestNewWeb3TokenFromBearer(t *testing.T) {
	bearer64 := "Bearer eyJzaWduZXIiOiJiaXRzb25nMXpxNjhkeDQyM2Zydjh5c3M1c2txeWc1dW05N3ZwZWZlZmUyZW5xIiwicGF5bG9hZCI6eyJkb21haW4iOiJ0ZXN0LmNvbSIsImV4cGlyZV9hdCI6MTY1NzI5MzUwNX0sInB1Yl9rZXkiOnsidHlwZSI6InRlbmRlcm1pbnQvUHViS2V5U2VjcDI1NmsxIiwidmFsdWUiOiJBdWpCUGtqMWNMVEZOeTVPbXVQNlBDRzd0dGtMSHhDaXdVQmdHNWdOYUpaTSJ9LCJzaWduYXR1cmUiOiJrY1FpVVE5NXFHdy9UU0pCWjJqS2JScXcyb2h3blhmTms4WFB0ZVQrU0FaVzNORzNxODE5OGNaeDNYZG9PY2EvKzlodWcxZHdkK2cwdnFaZnF1UHl1UT09In0="

	_, err := NewWeb3TokenFromBearer(bearer64)
	require.Nil(t, err)
}

func TestExpired(t *testing.T) {

	fmt.Println(time.Now().Add(time.Minute * 10).Unix())

	now := time.Now().Unix()
	fmt.Println(now)

	maxTime := time.Now().Add(MaxTimeLength)
	fmt.Println(maxTime)

	parsed := time.Unix(1111111111111, 0)
	fmt.Println(parsed)

	if maxTime.Before(parsed) {
		fmt.Println("invalid session")
	}
}
