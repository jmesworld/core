package interchaintest

import (
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	helpers "github.com/JMESWorld/core/tests/interchaintest/helpers"
)

// TestJunoDrip ensures the drip module properly distributes tokens from whitelisted accounts.
func TestJunoDrip(t *testing.T) {
	t.Parallel()

	// Setup new pre determined user (from test_node.sh)
	mnemonic := "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry"
	addr := "jmes1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl"

	// Base setup
	newCfg := jmesConfig
	newCfg.ModifyGenesis = cosmos.ModifyGenesis(append(defaultGenesisKV, []cosmos.GenesisKV{
		{
			Key:   "app_state.drip.params.allowed_addresses",
			Value: []string{addr},
		},
	}...))

	chains := CreateChainWithCustomConfig(t, 1, 0, newCfg)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	jmes := chains[0].(*cosmos.CosmosChain)
	nativeDenom := jmes.Config().Denom

	// User
	user, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "default", mnemonic, int64(1000000_000_000), jmes)
	if err != nil {
		t.Fatal(err)
	}

	// New TF token to distributes
	tfDenom := helpers.CreateTokenFactoryDenom(t, ctx, jmes, user, "dripme", fmt.Sprintf("0%s", Denom))
	distributeAmt := uint64(1_000_000)
	helpers.MintTokenFactoryDenom(t, ctx, jmes, user, distributeAmt, tfDenom)
	if balance, err := jmes.GetBalance(ctx, user.FormattedAddress(), tfDenom); err != nil {
		t.Fatal(err)
	} else if uint64(balance) != distributeAmt {
		t.Fatalf("balance not %d, got %d", distributeAmt, balance)
	}

	// Stake some tokens
	vals := helpers.GetValidators(t, ctx, jmes)
	valoper := vals.Validators[0].OperatorAddress

	stakeAmt := 100000_000_000
	helpers.StakeTokens(t, ctx, jmes, user, valoper, fmt.Sprintf("%d%s", stakeAmt, nativeDenom))

	// Drip the TF Tokens to all stakers
	distribute := int64(1_000_000)
	helpers.DripTokens(t, ctx, jmes, user, fmt.Sprintf("%d%s", distribute, tfDenom))

	// Claim staking rewards to capture the drip
	helpers.ClaimStakingRewards(t, ctx, jmes, user, valoper)

	// Check balances has the TF Denom from the claim
	bals, _ := jmes.AllBalances(ctx, user.FormattedAddress())
	fmt.Println("balances", bals)

	found := false
	for _, bal := range bals {
		if bal.Denom == tfDenom {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("did not find drip token")
	}

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
