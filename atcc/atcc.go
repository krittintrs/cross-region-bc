package main

import (
	"log"
	"atcc/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)



func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating myCC chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting myCC chaincode: %v", err)
	}
}
