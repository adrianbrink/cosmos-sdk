package bank

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

func unknownRequest(prefix string, ty interface{}) sdk.Result {
	errMsg := prefix + reflect.TypeOf(ty).Name()
	return sdk.ErrUnknownRequest(errMsg).Result()
}

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ibc.MsgSend:
			switch p := msg.Payload.(type) {
			case PayloadCoins:
				return handlePayloadCoinsSend(ctx, k, p, msg.DestChain)
			default:
				return unknownRequest("Unrecognized ibc/bank payload type: ", p)
			}
		case ibc.MsgReceive:
			return k.ch.Receive(func(ctx sdk.Context, p ibc.Payload) (ibc.Payload, sdk.Result) {
				switch p := msg.Payload.(type) {
				case PayloadCoins:
					return handlePayloadCoinsReceive(ctx, k, p)
				default:
					return nil, unknownRequest("Unrecognized ibc/bank payload type: ", p)
				}

			}, ctx, msg)
		// case ibc.MsgRelay
		default:
			return unknownRequest("Unrecognized ibc/bank Msg type: ", msg)
		}
	}
}

func handlePayloadCoinsSend(ctx sdk.Context, k Keeper, p PayloadCoins, chainid string) sdk.Result {
	tags, err := k.sendCoins(ctx, p, chainid)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Tags: tags}
}

func handlePayloadCoinsReceive(ctx sdk.Context, k Keeper, p PayloadCoins) (ibc.Payload, sdk.Result) {
	tags, err := k.receiveCoins(ctx, p)
	if err != nil {
		return PayloadCoinsFail{p}, err.Result()
	}
	return nil, sdk.Result{Tags: tags}
}
