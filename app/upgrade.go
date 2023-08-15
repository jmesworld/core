package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeHandler h for software upgrade proposal
type UpgradeHandler struct {
	*JmesApp
}

// NewUpgradeHandler return new instance of UpgradeHandler
func NewUpgradeHandler(app *JmesApp) UpgradeHandler {
	return UpgradeHandler{app}
}

func (h UpgradeHandler) CreateUpgradeHandler() upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return h.mm.RunMigrations(ctx, h.configurator, vm)
	}
}

func (h UpgradeHandler) vestingScheduleUpdateHandler(ctx sdk.Context, genesisTime int64, unlockAmountMap map[string]int64) error {
	bondDenom := h.StakingKeeper.BondDenom(ctx)
	for address, unlockAmount := range unlockAmountMap {
		accAddr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}

		// Required tokens are already allocated at genesis
		// but only vesting schedules are not properly set.
		// Thus, need to unlock tokens from vesting tokens
		account := h.AccountKeeper.GetAccount(ctx, accAddr)
		vestingAccount := account.(*vestingtypes.PeriodicVestingAccount)
		vestingAccount.OriginalVesting = vestingAccount.OriginalVesting.Sub(
			sdk.NewCoin(bondDenom, sdk.NewInt(unlockAmount)),
		)

		// Track delegation - decrease delegated vesting
		// and increase delegated free amount
		originalVesting := vestingAccount.OriginalVesting.AmountOf(bondDenom)
		delegatedVesting := vestingAccount.DelegatedVesting.AmountOf(bondDenom)
		delegatedFree := vestingAccount.DelegatedFree.AmountOf(bondDenom)
		delegatedAmount := delegatedFree.Add(delegatedVesting)
		if delegatedVesting.GT(originalVesting) {
			vestingAccount.DelegatedVesting = sdk.NewCoins(sdk.NewCoin(bondDenom, originalVesting))
			vestingAccount.DelegatedFree = sdk.NewCoins(sdk.NewCoin(bondDenom, delegatedAmount.Sub(originalVesting)))
		}

		// 2 year vesting with 6 month cliff
		vestingAccount.StartTime = genesisTime + 60*60*24*30*6
		vestingAccount.VestingPeriods = vestingtypes.Periods{
			{
				Length: 60 * 60 * 24 * 365 * 2,
				Amount: vestingAccount.OriginalVesting,
			},
		}

		// update account
		h.AccountKeeper.SetAccount(ctx, vestingAccount)
	}

	return nil
}
