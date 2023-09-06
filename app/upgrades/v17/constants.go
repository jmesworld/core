package v16

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/jmesworld/core/v17/app/upgrades"
	driptypes "github.com/jmesworld/core/v17/x/drip/types"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v17"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV17UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			driptypes.ModuleName,
		},
	},
}
