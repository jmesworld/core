#!/bin/bash
# microtick and bitcanna contributed significantly here.
# Pebbledb state sync script.
set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Initialize chain.
jmesd init test

# Get Genesis
wget https://download.dimi.sh/jmes-phoenix2-genesis.tar.gz
tar -xvf jmes-phoenix2-genesis.tar.gz
mv jmes-phoenix2-genesis.json "$HOME/.jmes/config/genesis.json"




# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT="$(curl -s https://jmes-rpc.polkachu.com/block | jq -r .result.block.header.height)"
BLOCK_HEIGHT="$((LATEST_HEIGHT-INTERVAL))"
TRUST_HASH="$(curl -s "https://jmes-rpc.polkachu.com/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)"

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export JUNOD_STATESYNC_ENABLE=true
export JUNOD_P2P_MAX_NUM_OUTBOUND_PEERS=200
export JUNOD_STATESYNC_RPC_SERVERS="https://rpc-jmes-ia.notional.ventures:443,https://jmes-rpc.polkachu.com:443"
export JUNOD_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export JUNOD_STATESYNC_TRUST_HASH=$TRUST_HASH

# Fetch and set list of seeds from chain registry.
JUNOD_P2P_SEEDS="$(curl -s https://raw.githubusercontent.com/cosmos/chain-registry/master/jmes/chain.json | jq -r '[foreach .peers.seeds[] as $item (""; "\($item.id)@\($item.address)")] | join(",")')"
export JUNOD_P2P_SEEDS

# Start chain.
jmesd start --x-crisis-skip-assert-invariants
