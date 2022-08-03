#!/bin/bash 

__dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. ${__dir}/vars/variables.sh

lavad tx gov submit-proposal spec-add ./cookbook/spec_add_ethereum.json,./cookbook/spec_add_osmosis.json,./cookbook/spec_add_fantom.json --from alice --gas-adjustment "1.5" --gas "auto" -y
lavad tx gov vote 1 yes -y --from alice
sleep 4

lavad tx pairing stake-client "ETH1" 200000ulava 1 -y --from user1
lavad tx pairing stake-client "COS3" 200000ulava 1 -y --from user2
lavad tx pairing stake-client "FTM250" 200000ulava 1 -y --from user3
lavad tx pairing stake-client "COS4" 200000ulava 1 -y --from user2

#Ethereum providers
lavad tx pairing stake-provider "ETH1" 2010ulava "127.0.0.1:2221,jsonrpc,1" 1 -y --from servicer1
lavad tx pairing stake-provider "ETH1" 2000ulava "127.0.0.1:2222,jsonrpc,1" 1 -y --from servicer2
lavad tx pairing stake-provider "ETH1" 2050ulava "127.0.0.1:2223,jsonrpc,1" 1 -y --from servicer3
lavad tx pairing stake-provider "ETH1" 2020ulava "127.0.0.1:2224,jsonrpc,1" 1 -y --from servicer4
lavad tx pairing stake-provider "ETH1" 2030ulava "127.0.0.1:2225,jsonrpc,1" 1 -y --from servicer5

#Osmosis providers
lavad tx pairing stake-provider "COS3" 2010ulava "127.0.0.1:2241,tendermintrpc,1 127.0.0.1:2231,rest,1" 1 -y --from servicer1
lavad tx pairing stake-provider "COS3" 2000ulava "127.0.0.1:2242,tendermintrpc,1 127.0.0.1:2232,rest,1" 1 -y --from servicer2
lavad tx pairing stake-provider "COS3" 2050ulava "127.0.0.1:2243,tendermintrpc,1 127.0.0.1:2233,rest,1" 1 -y --from servicer3

#Fantom providers
lavad tx pairing stake-provider "FTM250" 2010ulava "127.0.0.1:2251,jsonrpc,1" 1 -y --from servicer1
lavad tx pairing stake-provider "FTM250" 2000ulava "127.0.0.1:2252,jsonrpc,1" 1 -y --from servicer2
lavad tx pairing stake-provider "FTM250" 2050ulava "127.0.0.1:2253,jsonrpc,1" 1 -y --from servicer3
lavad tx pairing stake-provider "FTM250" 2020ulava "127.0.0.1:2254,jsonrpc,1" 1 -y --from servicer4
lavad tx pairing stake-provider "FTM250" 2030ulava "127.0.0.1:2255,jsonrpc,1" 1 -y --from servicer5


#Osmosis testnet providers
lavad tx pairing stake-provider "COS4" 2010ulava "127.0.0.1:4241,tendermintrpc,1 127.0.0.1:4231,rest,1" 1 -y --from servicer1
lavad tx pairing stake-provider "COS4" 2000ulava "127.0.0.1:4242,tendermintrpc,1 127.0.0.1:4232,rest,1" 1 -y --from servicer2
lavad tx pairing stake-provider "COS4" 2050ulava "127.0.0.1:4243,tendermintrpc,1 127.0.0.1:4233,rest,1" 1 -y --from servicer3


echo "---------------Queries------------------"
lavad query pairing providers "ETH1"
lavad query pairing clients "ETH1"


. ${__dir}/setup_providers.sh