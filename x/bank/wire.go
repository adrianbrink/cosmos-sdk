package bank

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/Send", nil)
	cdc.RegisterConcrete(MsgIssue{}, "cosmos-sdk/Issue", nil)
	cdc.RegisterConcrete(MsgIBCSend{}, "cosmos-sdk/IBCSend", nil)

	cdc.RegisterConcrete(PayloadSend{}, "cosmos-sdk/ibc/Send", nil)
	cdc.RegisterConcrete(ReceiptSendFail{}, "cosmos-sdk/ibc/SendFail", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
