package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/lib"
	"github.com/cosmos/cosmos-sdk/wire"
)

// ------------------------------------------
// Type Definitions

// Keeper manages connection between chains
type Keeper struct {
	key sdk.StoreKey
	cdc *wire.Codec

	codespace sdk.CodespaceType
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,

		codespace: codespace,
	}
}

// -----------------------------------------
// Store Accessors

func CommitHeightKey(srcChain string) []byte {
	return append([]byte{0x00}, []byte(srcChain)...)
}

func commitHeight(store sdk.KVStore, cdc *wire.Codec, srcChain string) (res uint64, ok bool) {
	bz := store.Get(CommitHeightKey(srcChain))
	if bz == nil {
		return res, false
	}
	cdc.MustUnmarshalBinary(bz, &res)
	return res, true
}

func CommitListPrefix(srcChain string) []byte {
	return append([]byte{0x01}, []byte(srcChain)...)
}

func commitList(store sdk.KVStore, cdc *wire.Codec, srcChain string) lib.List {
	return lib.NewList(cdc, store.Prefix(CommitListPrefix(srcChain)), nil)
}
