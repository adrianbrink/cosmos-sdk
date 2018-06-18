package store

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClientStore struct {
	ctx context.CoreContext
	storeName
}
