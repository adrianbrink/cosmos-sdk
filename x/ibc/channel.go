package ibc

import (
	"github.com/tendermint/go-amino"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/lib"
	"github.com/cosmos/cosmos-sdk/wire"
)

// ------------------------------------------
// Type Definitions

type Channel struct {
	k   Keeper
	key sdk.KVStoreGetter
}

func (k Keeper) Channel(key sdk.KVStoreGetter) Channel {
	return Channel{
		k:   k,
		key: key,
	}
}

type DatagramType byte

const (
	PacketType = DatagramType(iota)
	ReceiptType
)

type Header struct {
	SrcChain  string
	DestChain string
}

func (h Header) InverseDirection() Header {
	return Header{
		SrcChain:  h.DestChain,
		DestChain: h.SrcChain,
	}
}

type Payload interface {
	Type() string
	ValidateBasic() sdk.Error
	GetSigners() []sdk.AccAddress
	DatagramType() DatagramType
}

type Datagram struct {
	Header
	// Should we unexport Payload to possible modification from the modules?
	Payload
}

type Proof struct {
	Height   uint64
	Sequence uint64
}

// -------------------------------------------
// Store Accessors

func OutgoingQueuePrefix(ty DatagramType, destChain string) []byte {
	return append(append([]byte{0x00}, byte(ty)), []byte(destChain)...)
}

func outgoingQueue(store sdk.KVStore, cdc *wire.Codec, ty DatagramType, destChain string) lib.Linear {
	return lib.NewLinear(cdc, store.Prefix(OutgoingQueuePrefix(ty, destChain)), nil)
}

func IncomingSequenceKey(ty DatagramType, srcChain string) []byte {
	return append(append([]byte{0x01}, byte(ty)), []byte(srcChain)...)
}

func getIncomingSequence(store sdk.KVStore, ty DatagramType, srcChain string) (res uint64) {
	bz := store.Get(IncomingSequenceKey(ty, srcChain))
	amino.MustUnmarshalBinary(bz, &res)
	return
}

func setIncomingSequence(store sdk.KVStore, ty DatagramType, srcChain string, seq uint64) {
	bz := amino.MustMarshalBinary(seq)
	store.Set(IncomingSequenceKey(ty, srcChain), bz)
}
