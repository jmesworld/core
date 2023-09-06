#!/bin/sh

export JUNOD_NODE="http://localhost:26657"
CHAIN_A_ARGS="--from jmes1 --keyring-backend test --chain-id local-1 --home $HOME/.jmes1/ --node http://localhost:26657 --yes"

# jmesd q ibc channel channels

# Send from local-1 to local-2 via the relayer
jmesd tx ibc-transfer transfer transfer channel-0 jmes1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl 9ujmes $CHAIN_A_ARGS --packet-timeout-height 0-0

sleep 6

# check the query on the other chain to ensure it went through
jmesd q bank balances jmes1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl --chain-id local-2 --node http://localhost:36657