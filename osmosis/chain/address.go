// https://github.com/strangelove-ventures/lens/blob/main/client/address.go

package chain

import (
	"fmt"
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
func (c *Client) EncodeBech32AccPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", c.config.AccountPrefix, "pub"), addr)
}
func (c *Client) EncodeBech32ValAddr(addr sdk.ValAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", c.config.AccountPrefix, "valoper"), addr)
}
func (c *Client) MustEncodeValAddr(addr sdk.ValAddress) string {
	enc, err := c.EncodeBech32ValAddr(addr)
	if err != nil {
		panic(err)
	}
	return enc
}
func (c *Client) EncodeBech32ValPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", c.config.AccountPrefix, "valoperpub"), addr)
}
func (c *Client) EncodeBech32ConsAddr(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", c.config.AccountPrefix, "valcons"), addr)
}
func (c *Client) EncodeBech32ConsPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", c.config.AccountPrefix, "valconspub"), addr)
}

func (c *Client) DecodeBech32AccAddr(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, c.config.AccountPrefix)
}
func (c *Client) DecodeBech32AccPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", c.config.AccountPrefix, "pub"))
}
func (c *Client) DecodeBech32ValAddr(addr string) (sdk.ValAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", c.config.AccountPrefix, "valoper"))
}
func (c *Client) DecodeBech32ValPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", c.config.AccountPrefix, "valoperpub"))
}
func (c *Client) DecodeBech32ConsAddr(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", c.config.AccountPrefix, "valcons"))
}
func (c *Client) DecodeBech32ConsPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", c.config.AccountPrefix, "valconspub"))
}
