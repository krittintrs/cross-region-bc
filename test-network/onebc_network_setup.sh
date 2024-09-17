#!/bin/bash

echo "[ ONE BC NETWORK SETUP ]"

cd ../test-network
./network.sh down
./network.sh up createChannel -c mychannel -ca 

./network.sh deployCC -ccn regionalCC1 -ccp ../crosschain/regional1 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode queryinstalled