package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PayloadCoins struct {
	SrcAddr  sdk.AccAddress `json:"src-addr"`
	DestAddr sdk.AccAddress `json:"dest-addr"`
	Coins    sdk.Coins      `json:"coins"`
}

func (p PayloadCoins) Type() string {
	return "ibc/bank"
}

func (p PayloadCoins) ValidateBasic() sdk.Error {
	if !p.Coins.IsValid() {
		return sdk.ErrInvalidCoins(p.Coins.String())
	}
	if !p.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(p.Coins.String())
	}
	return nil
}

type PayloadCoinsFail struct {
	PayloadCoins
}
