package ibc

import (
	"reflect"

	"github.com/tendermint/tendermint/lite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgOpenConnection:
			return handleMsgOpenConnection(ctx, k, msg)
		case MsgUpdateConnection:
			return handleMsgUpdateConnection(ctx, k, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgOpenConnection(ctx sdk.Context, k Keeper, msg MsgOpenConnection) sdk.Result {
	store := ctx.KVStore(k.key)

	_, established := commitHeight(store, k.cdc, msg.SrcChain)
	if established {
		return ErrConnectionAlreadyEstablished(k.codespace).Result()
	}

	height := uint64(msg.ROT.Height())
	commits := commitList(store, k.cdc, msg.SrcChain)
	commits.Set(height, msg.ROT)

	return sdk.Result{}
}

func handleMsgUpdateConnection(ctx sdk.Context, k Keeper, msg MsgUpdateConnection) sdk.Result {
	store := ctx.KVStore(k.key)
	lastheight, established := commitHeight(store, k.cdc, msg.SrcChain)
	if !established {
		return ErrConnectionNotEstablished(k.codespace).Result()
	}

	commits := commitList(store, k.cdc, msg.SrcChain)
	var lastcommit lite.Commit
	err := commits.Get(lastheight, lastcommit)
	if err != nil {
		panic(err)
	}

	// TODO: add lc verificatioon
	/*
		cert := lite.NewDynamicCertifier(msg.SrcChain, commit.Validators, height)
		if err := cert.Update(msg.Commit); err != nil {
			return ErrUpdateCommitFailed(k.codespace, err).Result()
		}

		k.setCommit(ctx, msg.SrcChain, msg.Commit.Height(), msg.Commit)
	*/
	height := uint64(msg.Commit.Commit.Height())
	if height < lastheight {
		return ErrInvalidHeight(k.codespace).Result()
	}
	commits.Set(height, msg.Commit)
	return sdk.Result{}
}
