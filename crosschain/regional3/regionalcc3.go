package main

import (
	"log"

	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an GlobalAsset
type SmartContract struct {
	contractapi.Contract
}

type RegionalAsset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

func generateRegionalAssets(numAssets int) []RegionalAsset {
	var assets []RegionalAsset
	for i := 1; i <= numAssets; i++ {
		asset := RegionalAsset{
			ID:             fmt.Sprintf("pc%d", i),
			Color:          "yellow",       // Example color
			Size:           300,            // Example size
			Owner:          "PATIENT RG 3", // Example owner
			AppraisedValue: 30,             // Example appraised value
		}
		assets = append(assets, asset)
	}
	fmt.Printf("Successfully generated %d records\n", numAssets)
	return assets
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, numRows int) error {

	assets := generateRegionalAssets(numRows)

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			fmt.Printf("INIT: MARSHAL ERR!\n")
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			fmt.Printf("INIT: PUT STATE ERR!\n")
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	// fmt.Printf("Initial data added!\n")
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := RegionalAsset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*RegionalAsset, error) {
	startTime := time.Now()
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("FAIL TO READ!! Time taken for query for id (%s): %s\n", id, duration)

		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		duration := time.Since(startTime)
		fmt.Printf("NOT FOUND!! Time taken for query for id (%s): %s\n", id, duration)

		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset RegionalAsset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("UNMARSHAL ERR!! Time taken for query for id (%s): %s\n", id, duration)

		return nil, err
	}

	duration := time.Since(startTime)
	fmt.Printf("Time taken for query for id (%s): %s\n", id, duration)

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := RegionalAsset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*RegionalAsset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*RegionalAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset RegionalAsset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating regionalCC3 chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting regionalCC3 chaincode: %v", err)
	}
}
