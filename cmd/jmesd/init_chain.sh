#!/bin/bash

# Prompt user to clean ~/.jmes/
read -p "Do you want to clean ~/.jmes/ directory? (y/n): " clean_jmes

if [ "$clean_jmes" == "y" ]; then
    echo "Cleaning ~/.jmes/ directory..."
    rm -rf ~/.jmes/
    echo "Cleaned ~/.jmes/ directory."
fi

default_chain_id="jmes-namito"

read -p "Enter the chain-id: (default: $default_chain_id) " chain_id # jmes-namito

if [ -z "$chain_id" ]; then
    chain_id=$default_chain_id
fi

# Prompt user for moniker
read -p "Enter the moniker for the key: " moniker

# Prompt user to add a new key
read -p "Do you want to create a new key? It will overwrite prior created key with same name (y/n): " add_key

if [ "$add_key" == "y" ]; then
    # Add a new key
    echo "Adding a new key with moniker '$moniker'..."
    ./jmesd keys add "$moniker"
    if [ $? -eq 0 ]; then
        echo "New key added successfully."
    else
        echo "Error adding a new key."
        exit 1
    fi
fi

# Define the rest of the commands
commands=(
    "./jmesd init --chain-id=$chain_id $moniker"
    "./jmesd genesis add-genesis-account $moniker 0ujmes"
    "./jmesd genesis add-genesis-account jmes1mr4gj98n2cy2fnnme87gwctv34axsc63vjhyn9 45000000000000ujmes --vesting-amount 45000000000000ujmes --vesting-unlock-percentage 0.045"

    "./jmesd genesis add-genesis-account jmes1nm9rsr3yeuvvsdv3w2r7kqfjfa90tjwwd775rt 17000000000000ujmes --vesting-amount 17000000000000ujmes --vesting-unlock-percentage 0.017"
    "./jmesd genesis add-genesis-account jmes1v6nnc8f9jjpw309xm6zhj968vscgvld04gmhq0 17000000000000ujmes --vesting-amount 17000000000000ujmes --vesting-unlock-percentage 0.017"

    "./jmesd genesis add-genesis-account jmes12h9pe8v2pmqzec0lx8z5d6wsxkuell02nyzt70 16000000000000ujmes --vesting-amount 16000000000000ujmes --vesting-unlock-percentage 0.016"
    "./jmesd genesis add-genesis-account jmes1jlwkusyhyr9l96mgkv25uytr6wcvqn3whs6kez 5000000000000ujmes --vesting-amount 5000000000000ujmes --vesting-unlock-percentage 0.005"

    "./jmesd genesis add-genesis-account jmes1afze5vvplveqkha063q7fu46dcvylu7nzdax99 1000000000000ujmes --vesting-amount 1000000000000ujmes --vesting-unlock-percentage 0.001"
    "./jmesd genesis add-genesis-account jmes1uyauwmjasmq3hqcvpjmlxm72l2yrfjjt5v0nww 1000000000000ujmes --vesting-amount 1000000000000ujmes --vesting-unlock-percentage 0.001"
    "./jmesd genesis add-genesis-account jmes1v6ht8shgvs645lvrjejqaknw6d2cgcsnhk98nl 1000000000000ujmes --vesting-amount 1000000000000ujmes --vesting-unlock-percentage 0.001"

    "./jmesd genesis add-genesis-account jmes1pmcm6ag8hn7y6009q5e3q4dga9268epgxm2r6y 1000000000000ujmes --vesting-amount 1000000000000ujmes --vesting-unlock-percentage 0.001"

    "./jmesd genesis gentx $moniker 0ujmes --chain-id=$chain_id"
    "./jmesd genesis collect-gentxs"
)

# Execute the commands one by one
for cmd in "${commands[@]}"; do
    echo "Running: $cmd"
    eval "$cmd"
    if [ $? -eq 0 ]; then
        echo "Command successfully executed: $cmd"
    else
        echo "Error executing command: $cmd"
        exit 1
    fi
done

echo "Initialization complete."