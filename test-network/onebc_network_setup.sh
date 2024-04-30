#!/bin/bash

echo "[ ONE BC NETWORK SETUP ]"

cd ../test-network
./network.sh down
./network.sh up createChannel -c mychannel -ca 

./network.sh deployCC -ccn regionalCC1 -ccp ../crosschain/regional1 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

peer lifecycle chaincode queryinstalled