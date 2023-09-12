package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jmesworld/core/v2/app/config"
)

func RegisterDenomsConfig() error {
	// sdk.RegisterDenom(config.Jmes, sdk.OneDec())
	// sdk.RegisterDenom(config.MilliJmes, sdk.NewDecWithPrec(1, 3))
	err := sdk.RegisterDenom(config.MicroJmes, sdk.NewDecWithPrec(1, 6))
	if err != nil {
		return err
	}
	// sdk.RegisterDenom(config.NanoJmes, sdk.NewDecWithPrec(1, 9))

	return nil
}
