package feeshare_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/require"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	jmesapp "github.com/jmesworld/core/v17/app"
	"github.com/jmesworld/core/v17/x/mint/types"
)

// returns context and an app with updated mint keeper
func CreateTestApp(t *testing.T, isCheckTx bool) (*jmesapp.App, sdk.Context) {
	app := Setup(t, isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: "testing",
	})
	if err := app.AppKeepers.MintKeeper.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}
	app.AppKeepers.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func Setup(t *testing.T, isCheckTx bool) *jmesapp.App {
	app, genesisState := GenApp(t, !isCheckTx)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators: []abci.ValidatorUpdate{},
				// ConsensusParams: &tmproto.ConsensusParams{},
				ConsensusParams: jmesapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
				ChainId:         "testing",
			},
		)
	}

	return app
}

func GenApp(t *testing.T, withGenesis bool, opts ...wasmkeeper.Option) (*jmesapp.App, jmesapp.GenesisState) {
	db := dbm.NewMemDB()
	nodeHome := t.TempDir()
	snapshotDir := filepath.Join(nodeHome, "data", "snapshots")

	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	require.NoError(t, err)
	t.Cleanup(func() { snapshotDB.Close() })
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	require.NoError(t, err)

	app := jmesapp.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		wasmtypes.EnableAllProposals,
		simtestutil.EmptyAppOptions{},
		opts,
		bam.SetChainID("testing"),
		bam.SetSnapshot(snapshotStore, snapshottypes.SnapshotOptions{KeepRecent: 2}),
	)

	if withGenesis {
		return app, jmesapp.NewDefaultGenesisState(app.AppCodec())
	}

	return app, jmesapp.GenesisState{}
}
