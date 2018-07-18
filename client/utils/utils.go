package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
)

// GetPassphraseFromStdin gets the passphrase from STDIN. An error is returned
// if getting input fails.
func GetPassphraseFromStdin(name string) (string, error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Passphrase to sign with '%s':", name)

	return client.GetPassword(prompt, buf)
}

// SendTx implements a handler that facilitates sending a series of messages in
// a signed transaction given a TxContext and a QueryContext. It ensures that
// the account has a proper number and sequence set. In addition, it builds and
// signs a transaction with the supplied messages. Finally, it broadcasts the
// signed transaction to a node.
func SendTx(txCtx authctx.TxContext, queryCtx context.QueryContext, from []byte, name string, msgs []sdk.Msg) error {
	if txCtx.AccountNumber == 0 {
		accNum, err := queryCtx.GetAccountNumber(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithAccountNumber(accNum)
	}

	if txCtx.Sequence == 0 {
		accSeq, err := queryCtx.GetAccountSequence(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithSequence(accSeq)
	}

	keyInfo, err := keys.GetKeyInfo(name)
	if err != nil {
		return err
	}

	var passphrase string

	// we only need a passphrase for locally stored keys
	if keyInfo.GetType() == cryptokeys.TypeLocal {
		passphrase, err = GetPassphraseFromStdin(name)
		if err != nil {
			return fmt.Errorf("Error reading passphrase: %v", err)
		}
	}

	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(name, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to Tendermint
	if err := queryCtx.EnsureBroadcastTx(txBytes); err != nil {
		return err
	}

	return nil
}
