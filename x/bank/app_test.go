package bank

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

type appTestCase struct {
	addr     sdk.AccAddress
	coins    sdk.Coins
	simBlock bool
	expPass  bool
	msgs     []sdk.Msg
	accNums  []int64
	accSeqs  []int64
	privKeys []crypto.PrivKey
}

var (
	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = crypto.GenPrivKeyEd25519()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(crypto.GenPrivKeyEd25519().PubKey().Address())
	priv4 = crypto.GenPrivKeyEd25519()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

	coins     = sdk.Coins{sdk.NewCoin("foocoin", 10)}
	halfCoins = sdk.Coins{sdk.NewCoin("foocoin", 5)}
	manyCoins = sdk.Coins{sdk.NewCoin("foocoin", 1), sdk.NewCoin("barcoin", 1)}
	freeFee   = auth.NewStdFee(100000, sdk.Coins{sdk.NewCoin("foocoin", 0)}...)

	sendMsg1 = MsgSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}
	sendMsg2 = MsgSend{
		Inputs: []Input{NewInput(addr1, coins)},
		Outputs: []Output{
			NewOutput(addr2, halfCoins),
			NewOutput(addr3, halfCoins),
		},
	}
	sendMsg3 = MsgSend{
		Inputs: []Input{
			NewInput(addr1, coins),
			NewInput(addr4, coins),
		},
		Outputs: []Output{
			NewOutput(addr2, coins),
			NewOutput(addr3, coins),
		},
	}
	sendMsg4 = MsgSend{
		Inputs: []Input{
			NewInput(addr2, coins),
		},
		Outputs: []Output{
			NewOutput(addr1, coins),
		},
	}
	sendMsg5 = MsgSend{
		Inputs: []Input{
			NewInput(addr1, manyCoins),
		},
		Outputs: []Output{
			NewOutput(addr2, manyCoins),
		},
	}
)

// initialize the mock application for this module
func getMockApp(t *testing.T) *mock.App {
	mapp, err := getBenchmarkMockApp()
	require.NoError(t, err)
	return mapp
}

func TestBankWithRandomMessages(t *testing.T) {
	mapp := getMockApp(t)
	setup := func(r *rand.Rand, keys []crypto.PrivKey) {
		return
	}

	mapp.RandomizedTesting(
		t,
		[]mock.TestAndRunTx{TestAndRunSingleInputMsgSend},
		[]mock.RandSetup{setup},
		[]mock.Invariant{ModuleInvariants},
		100, 30, 30,
	)
}

func TestMsgSendWithAccounts(t *testing.T) {
	mapp := getMockApp(t)
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 67)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc})

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	res1 := mapp.AccountMapper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	testCases := []appTestCase{
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg1},
			accNums:  []int64{0},
			accSeqs:  []int64{0},
			expPass:  true,
			privKeys: []crypto.PrivKey{priv1},
		},
		{
			addr:  addr1,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 57)},
		},
		{
			addr:  addr2,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 10)},
		},
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg1, sendMsg2},
			accNums:  []int64{0},
			accSeqs:  []int64{0},
			expPass:  false,
			privKeys: []crypto.PrivKey{priv1},
		},
	}

	for _, tc := range testCases {
		if tc.simBlock {
			mock.SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expPass, tc.privKeys...)
		} else {
			mock.CheckBalance(t, mapp, tc.addr, tc.coins)
		}
	}

	// bumping the tx nonce number without resigning should be an auth error
	mapp.BeginBlock(abci.RequestBeginBlock{})

	tx := mock.GenTx([]sdk.Msg{sendMsg1}, []int64{0}, []int64{0}, priv1)
	tx.Signatures[0].Sequence = 1

	res := mapp.Deliver(tx)
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceRoot, sdk.CodeUnauthorized), res.Code, res.Log)

	// resigning the tx with the bumped sequence should work
	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{sendMsg1, sendMsg2}, []int64{0}, []int64{1}, true, priv1)
}

func TestMsgSendMultipleOut(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc1, acc2})

	testCases := []appTestCase{
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg2},
			accNums:  []int64{0},
			accSeqs:  []int64{0},
			expPass:  true,
			privKeys: []crypto.PrivKey{priv1},
		},
		{
			addr:  addr1,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 32)},
		},
		{
			addr:  addr2,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 47)},
		},
		{
			addr:  addr3,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 5)},
		},
	}

	for _, tc := range testCases {
		if tc.simBlock {
			mock.SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expPass, tc.privKeys...)
		} else {
			mock.CheckBalance(t, mapp, tc.addr, tc.coins)
		}
	}
}

func TestSengMsgMultipleInOut(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}
	acc4 := &auth.BaseAccount{
		Address: addr4,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc1, acc2, acc4})

	testCases := []appTestCase{
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg3},
			accNums:  []int64{0, 2},
			accSeqs:  []int64{0, 0},
			expPass:  true,
			privKeys: []crypto.PrivKey{priv1, priv4},
		},
		{
			addr:  addr1,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 32)},
		},
		{
			addr:  addr4,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 32)},
		},
		{
			addr:  addr2,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 52)},
		},
		{
			addr:  addr3,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 10)},
		},
	}

	for _, tc := range testCases {
		if tc.simBlock {
			mock.SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expPass, tc.privKeys...)
		} else {
			mock.CheckBalance(t, mapp, tc.addr, tc.coins)
		}
	}
}

func TestMsgSendDependent(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewCoin("foocoin", 42)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc1})

	testCases := []appTestCase{
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg1},
			accNums:  []int64{0},
			accSeqs:  []int64{0},
			expPass:  true,
			privKeys: []crypto.PrivKey{priv1},
		},
		{
			addr:  addr1,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 32)},
		},
		{
			addr:  addr2,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 10)},
		},
		{
			simBlock: true,
			msgs:     []sdk.Msg{sendMsg4},
			accNums:  []int64{1},
			accSeqs:  []int64{0},
			expPass:  true,
			privKeys: []crypto.PrivKey{priv2},
		},
		{
			addr:  addr1,
			coins: sdk.Coins{sdk.NewCoin("foocoin", 42)},
		},
	}

	for _, tc := range testCases {
		if tc.simBlock {
			mock.SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expPass, tc.privKeys...)
		} else {
			mock.CheckBalance(t, mapp, tc.addr, tc.coins)
		}
	}
}
