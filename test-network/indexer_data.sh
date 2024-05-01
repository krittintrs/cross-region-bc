#!/bin/bash

echo "[ INDEXER INITIATE DATA ]"

# Check if the correct number of arguments is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <total_num>"
    exit 1
fi

# Set variables for the number of assets for each regionalCC
num_assets_regionalCC1=$(($1/10*5))
num_assets_regionalCC2=$(($1/10*3))
num_assets_regionalCC3=$(($1/10*2))

echo $num_assets_regionalCC1 $num_assets_regionalCC2 $num_assets_regionalCC3

# Invoke chaincode for each regionalCC with the specified number of assets
peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C mychannel \
    -n regionalCC1 \
    -c "{\"Args\":[\"InitLedger\", \"$num_assets_regionalCC1\"]}"

peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C mychannel \
    -n regionalCC2 \
    -c "{\"Args\":[\"InitLedger\", \"$num_assets_regionalCC2\"]}"

peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C mychannel \
    -n regionalCC3 \
    -c "{\"Args\":[\"InitLedger\", \"$num_assets_regionalCC3\"]}"