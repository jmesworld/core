package interchaintest

import (
	"fmt"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/stretchr/testify/require"

	cosmosproto "github.com/cosmos/gogoproto/proto"

	helpers "github.com/JMESWorld/core/tests/interchaintest/helpers"
)

// TestJunoSubmitUnityContract test to ensure the store code properly works on the contract
// - https://github.com/CosmosContracts/cw-unity-prop
func TestJunoUnityContractGovSubmit(t *testing.T) {
	t.Parallel()

	// Base setup
	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	jmes := chains[0].(*cosmos.CosmosChain)

	nativeDenom := jmes.Config().Denom

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(10000_000000), jmes, jmes)
	user := users[0]
	withdrawUser := users[1]
	withdrawAddr := withdrawUser.FormattedAddress()

	// Upload & init unity contract with no admin in test mode
	msg := fmt.Sprintf(`{"native_denom":"%s","withdraw_address":"%s","withdraw_delay_in_days":28}`, nativeDenom, withdrawAddr)
	_, contractAddr := helpers.SetupContract(t, ctx, jmes, user.KeyName(), "contracts/cw_unity_prop.wasm", msg)
	t.Log("testing Unity contractAddr", contractAddr)

	// send 2JUNO funds to the contract from user
	jmes.SendFunds(ctx, user.KeyName(), ibc.WalletAmount{Address: contractAddr, Denom: nativeDenom, Amount: 2000000})

	height, err := jmes.Height(ctx)
	require.NoError(t, err, "error fetching height")

	// Use cosmos messages, then build the proposal, and submit it.
	proposalMsgs := []cosmosproto.Message{
		&wasmtypes.MsgSudoContract{
			Authority: "jmes10d07y265gmmuvt4z0w9aw880jnsr700jvss730",
			Contract:  contractAddr,
			Msg:       []byte(fmt.Sprintf(`{"execute_send":{"amount":"1000000","recipient":"%s"}}`, withdrawAddr)),
		},
	}

	proposal, err := jmes.BuildProposal(proposalMsgs, "Prop Title", "description", "ipfs://CID", fmt.Sprintf(`1000000000%s`, nativeDenom))
	require.NoError(t, err, "error making proposal")

	txProp, err := jmes.SubmitProposal(ctx, user.KeyName(), proposal)
	t.Log("txProp", txProp)
	t.Log("err", err)

	proposalID := "1"

	err = jmes.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	// poll for proposal
	_, err = cosmos.PollForProposalStatus(ctx, jmes, height, height+50, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
