#!/bin/sh
#./jmesd keys add dev_masternode --recover && ./jmes --help
yes "boost top desk keen unusual scene entire belt cargo protect subject donor front dose narrow fruit square despair chat crush visual reform river decorate" | ./jmesd keys add devnode --recover || ./jmesd init --chain-id=jmes-888 devnode && ./jmesd add-genesis-account jmes1pmcm6ag8hn7y6009q5e3q4dga9268epgxm2r6y 1000000ujmes && yes "boost top desk keen unusual scene entire belt cargo protect subject donor front dose narrow fruit square despair chat crush visual reform river decorate" | ./jmesd gentx devnode 1000000ujmes --chain-id=jmes-888 && ./jmesd collect-gentxs
#
#- name: wallet
#  type: local
#  address: jmes1wrrjk4xurqk2ukaxr4hlrmsfgprxm075gnwacn
#  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A4o7xiulaT27XTv/I0EB8HYptmS5S3npaRGrVQiAYOP8"}'
#  mnemonic: ""


#**Important** write this mnemonic phrase in a safe place.
#It is the only way to recover your account if you ever forget your password.
#
#boost top desk keen unusual scene entire belt cargo protect subject donor front dose narrow fruit square despair chat crush visual reform river decorate


#➜  jmesd git:(cosmwasm) ✗ add-genesis-account $(terrad keys show <account-name> -a) 100000000uluna,1000usd

#./jmesd add-genesis-account jmes1wrrjk4xurqk2ukaxr4hlrmsfgprxm075gnwacn 1ujmes
#./jmesd gentx wallet 1ujmes --chain-id=jmes-888
#./jmesd collect-gentxs
