// https://github.com/strangelove-ventures/lens/blob/main/client/address.go

package chain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c *Client) EncodeBech32AccAddr(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(c.config.AccountPrefix, addr)
}
func (c *Client) MustEncodeAccAddr(addr sdk.AccAddress) string {
	enc, err := c.EncodeBech32AccAddr(addr)
	if err != nil {
		panic(err)
	}
	return enc
}
