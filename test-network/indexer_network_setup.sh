#!/bin/bash

echo "[ INDEXER BC NETWORK SETUP ]"

cd ../test-network
./network.sh down
./network.sh up createChannel -c mychannel -ca

./network.sh deployCC -ccn regionalCC1 -ccp ../crosschain/regional1 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"
./network.sh deployCC -ccn regionalCC2 -ccp ../crosschain/regional2 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"
./network.sh deployCC -ccn regionalCC3 -ccp ../crosschain/regional3 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

peer lifecycle chaincode queryinstalled