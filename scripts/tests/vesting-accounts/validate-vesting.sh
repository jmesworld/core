#!/bin/bash

echo ""
echo "############################"
echo "# Validate Vesting Account #"
echo "############################"
echo ""

BINARY=jmesd
CHAIN_DIR=$(pwd)/data
VESTING_FILE=$(pwd)/scripts/tests/vesting-accounts/vesting-periods.json
HIDDEN_VESTING_FILE=$(pwd)/scripts/tests/vesting-accounts/.vesting-periods.json

WALLET_3=$($BINARY keys show wallet3 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_4=$($BINARY keys show wallet4 -a --keyring-backend test --home $CHAIN_DIR/test-2)

echo "Checking the delegated vesting balance of wallet3 on chain test-2 to 90000000000 since 10000000000 is vesting"
WALLET_4_BALANCES=$($BINARY query bank balances $WALLET_4 --chain-id test-2 --node tcp://localhost:26657 -o json | jq -r '.balances[-1].amount')
if [[ "$WALLET_4_BALANCES" != "990000000000" ]]; then
    echo "Error: Expected a balance of 990000000000, got $WALLET_4_BALANCES"
    exit 1
fi

echo "Checking the vesting balance of wallet3 to be staked on chain test-2 to 10000000000"
WALLET_4_DELEGATIONS=$($BINARY query staking delegations $WALLET_4 --home $CHAIN_DIR/test-2 --node tcp://localhost:26657 -o json | jq -r '.delegation_responses[-1].balance.amount')
if [[ "$WALLET_4_DELEGATIONS" != "10000000000" ]]; then
    echo "Error: Expected a total staking of of 10000000000, got $WALLET_4_DELEGATIONS"
    exit 1
fi

echo "Creating a random vesting wallet on chain test-1"
CURRENT_DATE=$(date +%s)
$BINARY keys add wallet$CURRENT_DATE --home $CHAIN_DIR/test-1 --keyring-backend=test &> /dev/null
RANDOM_VESTING_WALLET=$($BINARY keys show wallet$CURRENT_DATE -a --keyring-backend test --home $CHAIN_DIR/test-1)

cp $VESTING_FILE $HIDDEN_VESTING_FILE
sed -i -e 's/"start_time": -1/"start_time": '$CURRENT_DATE'/g' $HIDDEN_VESTING_FILE

echo "Deploying a vesting account on chain test-1 with the address $RANDOM_VESTING_WALLET"
TX_HASH=$($BINARY tx vesting create-periodic-vesting-account $RANDOM_VESTING_WALLET $HIDDEN_VESTING_FILE --from $WALLET_3 --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657  --keyring-backend test -y -o json | jq -r '.txhash')
sleep 3
CREATE_VESTING_ACCOUNT_MSG_RES=$(jmesd query tx $TX_HASH -o josn --chain-id test-1 --home $CHAIN_DIR/test-1 --node tcp://localhost:16657 | jq -r '.logs[0].events[0].attributes[0].value')

if [[ "$CREATE_VESTING_ACCOUNT_MSG_RES" != "/cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount" ]]; then
    echo "Error: Expected a message type /cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount, got $CREATE_VESTING_ACCOUNT_MSG_RES"
    exit 1
fi

echo "Waiting 4 seconds for address $RANDOM_VESTING_WALLET to have spendable balance"

PREV_SPENDABLE_BALANCE="0"
CURRENT_SPENDABLE_BALANCE="0"
for i in {0..4};do
    sleep 3
    CURRENT_SPENDABLE_BALANCE=$(curl -s -X GET "http://localhost:1316/cosmos/bank/v1beta1/spendable_balances/$RANDOM_VESTING_WALLET" -H "accept: application/json" | jq -r '.balances[-1].amount')

    if [[ $PREV_SPENDABLE_BALANCE -ge $CURRENT_SPENDABLE_BALANCE ]]; then
        echo "Error: Expected a current block spendable balance ($CURRENT_SPENDABLE_BALANCE) greater than prev block ($PREV_SPENDABLE_BALANCE)"
        exit 1
    fi

    PREV_SPENDABLE_BALANCE=$CURRENT_SPENDABLE_BALANCE
    echo "Spendable balance of $RANDOM_VESTING_WALLET is $CURRENT_SPENDABLE_BALANCE on iteration $i"
done

rm -rf $HIDDEN_VESTING_FILE

echo ""
echo "#####################################"
echo "# SUCCESS: Validate Vesting Account #"
echo "#####################################"
echo ""
