package ante

// Used for the JMES ante handler so we can properly send 50% of fees to dAPP developers via fee share module

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	revtypes "github.com/jmesworld/core/v17/x/feeshare/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type FeeShareKeeper interface {
	GetParams(ctx sdk.Context) revtypes.Params
	GetFeeShare(ctx sdk.Context, contract sdk.Address) (revtypes.FeeShare, bool)
}
