package ibc

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/Send", nil)
	cdc.RegisterConcrete(MsgReceive{}, "cosmos-sdk/Receive", nil)
	cdc.RegisterConcrete(MsgCleanup{}, "cosmos-sdk/Cleanup", nil)
	cdc.RegisterConcrete(MsgOpenConnection{}, "cosmos-sdk/OpenConnection", nil)
	cdc.RegisterConcrete(MsgUpdateConnection{}, "cosmos-sdk/UpdateConnection", nil)

	cdc.RegisterInterface((*Payload)(nil), nil)
}
