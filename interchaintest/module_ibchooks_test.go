package interchaintest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	helpers "github.com/JMESWorld/core/tests/interchaintest/helpers"
)

// TestJunoIBCHooks ensures the ibc-hooks middleware from osmosis works.
func TestJunoIBCHooks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with JMES and jmes2
	numVals := 1
	numFullNodes := 0

	cfg2 := jmesConfig.Clone()
	cfg2.Name = "jmes-counterparty"
	cfg2.ChainID = "counterparty-2"

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "jmes",
			ChainConfig:   jmesConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "jmes",
			ChainConfig:   cfg2,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	jmes, jmes2 := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(jmes).
		AddChain(jmes2).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  jmes,
			Chain2:  jmes2,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, jmes, jmes2)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, jmes, jmes2)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	jmesUser, jmes2User := users[0], users[1]

	jmesUserAddr := jmesUser.FormattedAddress()
	// jmes2UserAddr := jmes2User.FormattedAddress()

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, jmes.Config().ChainID, jmes2.Config().ChainID)
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	_, contractAddr := helpers.SetupContract(t, ctx, jmes2, jmes2User.KeyName(), "contracts/ibchooks_counter.wasm", `{"count":0}`)

	// do an ibc transfer through the memo to the other chain.
	transfer := ibc.WalletAmount{
		Address: contractAddr,
		Denom:   jmes.Config().Denom,
		Amount:  int64(1),
	}

	memo := ibc.TransferOptions{
		Memo: fmt.Sprintf(`{"wasm":{"contract":"%s","msg":%s}}`, contractAddr, `{"increment":{}}`),
	}

	// Initial transfer. Account is created by the wasm execute is not so we must do this twice to properly set up
	transferTx, err := jmes.SendIBCTransfer(ctx, channel.ChannelID, jmesUser.KeyName(), transfer, memo)
	require.NoError(t, err)
	jmesHeight, err := jmes.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, jmes, jmesHeight-5, jmesHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Second time, this will make the counter == 1 since the account is now created.
	transferTx, err = jmes.SendIBCTransfer(ctx, channel.ChannelID, jmesUser.KeyName(), transfer, memo)
	require.NoError(t, err)
	jmesHeight, err = jmes.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, jmes, jmesHeight-5, jmesHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Get the address on the other chain's side
	addr := helpers.GetIBCHooksUserAddress(t, ctx, jmes, channel.ChannelID, jmesUserAddr)
	require.NotEmpty(t, addr)

	// Get funds on the receiving chain
	funds := helpers.GetIBCHookTotalFunds(t, ctx, jmes2, contractAddr, addr)
	require.Equal(t, int(1), len(funds.Data.TotalFunds))

	var ibcDenom string
	for _, coin := range funds.Data.TotalFunds {
		if strings.HasPrefix(coin.Denom, "ibc/") {
			ibcDenom = coin.Denom
			break
		}
	}
	require.NotEmpty(t, ibcDenom)

	// ensure the count also increased to 1 as expected.
	count := helpers.GetIBCHookCount(t, ctx, jmes2, contractAddr, addr)
	require.Equal(t, int64(1), count.Data.Count)
}
