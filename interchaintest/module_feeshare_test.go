package interchaintest

import (
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	helpers "github.com/JMESWorld/core/tests/interchaintest/helpers"
)

// TestJunoFeeShare ensures the feeshare module register and execute sharing functions work properly on smart contracts.
func TestJunoFeeShare(t *testing.T) {
	t.Parallel()

	// Base setup
	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	jmes := chains[0].(*cosmos.CosmosChain)

	nativeDenom := jmes.Config().Denom

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(10_000_000), jmes, jmes)
	user := users[0]
	feeRcvAddr := "jmes1v75wlkccpv7le3560zw32v2zjes5n0e7csr4qh"

	// Upload & init contract payment to another address
	_, contractAddr := helpers.SetupContract(t, ctx, jmes, user.KeyName(), "contracts/cw_template.wasm", `{"count":0}`)

	// register contract to a random address (since we are the creator, though not the admin)
	helpers.RegisterFeeShare(t, ctx, jmes, user, contractAddr, feeRcvAddr)
	if balance, err := jmes.GetBalance(ctx, feeRcvAddr, nativeDenom); err != nil {
		t.Fatal(err)
	} else if balance != 0 {
		t.Fatal("balance not 0")
	}

	// execute with a 10000 fee (so 5000 denom should be in the contract now with 50% feeshare default)
	helpers.ExecuteMsgWithFee(t, ctx, jmes, user, contractAddr, "", "10000"+nativeDenom, `{"increment":{}}`)

	// check balance of nativeDenom now
	if balance, err := jmes.GetBalance(ctx, feeRcvAddr, nativeDenom); err != nil {
		t.Fatal(err)
	} else if balance != 5000 {
		t.Fatal("balance not 5,000. it is ", balance, nativeDenom)
	}

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
