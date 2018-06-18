package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SendHandler func(Payload) sdk.Result

func (ch Channel) Send(h SendHandler, ctx sdk.Context, msg MsgSend) (result sdk.Result) {
	payload := msg.Payload

	// TODO: check validity of the payload; the module have to be permitted to send the payload
	result = h(msg.Payload)
	if !result.IsOK() {
		return
	}

	store := ch.key.KVStore(ctx)
	queue := outgoingQueue(store, ch.k.cdc, payload.DatagramType(), msg.DestChain)
	queue.Push(Datagram{
		Header: Header{
			SrcChain:  ctx.ChainID(),
			DestChain: msg.DestChain,
		},
		Payload: payload,
	})

	return
}

type ReceiveHandler func(sdk.Context, Payload) (Payload, sdk.Result)

func (ch Channel) Receive(h ReceiveHandler, ctx sdk.Context, msg MsgReceive) (res sdk.Result) {
	data := msg.Datagram
	prf := msg.Proof
	ty := data.Payload.DatagramType()
	srcChain := data.Header.SrcChain
	destChain := data.Header.DestChain

	if ctx.ChainID() != destChain {
		return ErrChainMismatch(ch.k.codespace).Result()
	}

	// TODO: verify merkle proof

	store := ch.key.KVStore(ctx)
	seq := getIncomingSequence(store, ty, srcChain)
	if seq != prf.Sequence {
		return ErrInvalidSequence(ch.k.codespace).Result()
	}
	setIncomingSequence(store, ty, srcChain, seq+1)

	switch ty {
	case PacketType:
		return ch.receivePacket(h, ctx, store, data)
	case ReceiptType:
		return ch.receiveReceipt(h, ctx, data)
	default:
		// Source zone sent invalid datagram, reorg needed
		return ErrUnknownDatagramType(ch.k.codespace).Result()
	}
}

func (ch Channel) receivePacket(h ReceiveHandler, ctx sdk.Context, store sdk.KVStore, data Datagram) (res sdk.Result) {
	// Packet handling can fail
	// If fails, reverts all execution done by DatagramHandler

	cctx, write := ctx.CacheContext()
	receipt, res := h(cctx, data.Payload)
	if receipt != nil {
		newdata := Datagram{
			Header:  data.Header.InverseDirection(),
			Payload: receipt,
		}

		queue := outgoingQueue(store, ch.k.cdc, ReceiptType, newdata.Header.DestChain)
		queue.Push(newdata)
	}
	if !res.IsOK() {
		return WrapResult(res)
	}
	write()

	return
}

func (ch Channel) receiveReceipt(h ReceiveHandler, ctx sdk.Context, data Datagram) (res sdk.Result) {
	// Receipt handling should not fail

	receipt, res := h(ctx, data.Payload)
	if !res.IsOK() {
		panic("IBC Receipt handler should not fail")
	}
	if receipt != nil {
		panic("IBC Receipt handler cannot return new receipt")
	}

	return
}

/*
func cleanup(store sdk.KVStore, cdc *wire.Codec, ty DatagramType, srcChain string) sdk.Result {
	queue := outgoingQueue(store, cdc, ty, srcChain)
}
*/
