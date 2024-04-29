package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an GlobalAsset
type SmartContract struct {
	contractapi.Contract
}

// GlobalAsset describes basic details of what makes up a simple asset
type GlobalAsset struct {
	HospitalID     string `json:"hospitalID"`
	RegionalCCName string `json:"regionalCCName"`
}

type RegionalAsset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	assets := []GlobalAsset{
		{HospitalID: "HP1", RegionalCCName: "regionalCC1"},
		{HospitalID: "HP2", RegionalCCName: "regionalCC2"},
		{HospitalID: "HP3", RegionalCCName: "regionalCC3"},
		{HospitalID: "HP4", RegionalCCName: "regionalCC4"},
		{HospitalID: "HP5", RegionalCCName: "regionalCC5"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			fmt.Printf("INIT: MARSHAL ERR!\n")
			return err
		}

		err = ctx.GetStub().PutState(asset.HospitalID, assetJSON)
		if err != nil {
			fmt.Printf("INIT: PUT STATE ERR!\n")
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	// fmt.Printf("Initial data added!\n")
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, hospitalID string, rccName string) error {
	exists, err := s.AssetExists(ctx, hospitalID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", hospitalID)
	}

	asset := GlobalAsset{
		HospitalID:     hospitalID,
		RegionalCCName: rccName,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hospitalID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*GlobalAsset, error) {
	startTime := time.Now()
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("FAIL TO READ!! Time taken for query for policyID (%s): %s\n", id, duration)

		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		duration := time.Since(startTime)
		fmt.Printf("NOT FOUND!! Time taken for query for id (%s): %s\n", id, duration)

		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset GlobalAsset
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

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadRegionalAsset(ctx contractapi.TransactionContextInterface, policyID string, hospitalID string) (*RegionalAsset, error) {
	startTime := time.Now()
	indexAssetJSON, err := ctx.GetStub().GetState(hospitalID)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("FAIL TO READ!! Time taken for query for hospitalID (%s): %s\n", hospitalID, duration)

		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if indexAssetJSON == nil {
		duration := time.Since(startTime)
		fmt.Printf("NOT FOUND!! Time taken for query for hospitalID (%s): %s\n", hospitalID, duration)

		return nil, fmt.Errorf("the asset hospitalID (%s) does not exist", hospitalID)
	}

	var indexAsset GlobalAsset
	err = json.Unmarshal(indexAssetJSON, &indexAsset)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("UNMARSHAL ERR!! Time taken for query for hospitalID (%s): %s\n", hospitalID, duration)

		return nil, err
	}
	
	duration := time.Since(startTime)
	fmt.Printf("Time before retrieve regional: %s\n", duration)

	asset, err := retrieveFromRegionalBC(ctx, indexAsset.RegionalCCName, policyID)
	if err != nil {
		duration := time.Since(startTime)
		fmt.Printf("Regional ERR!! Time taken for query for hospitalID (%s): %s\n", hospitalID, duration)

		return nil, err
	}

	duration = time.Since(startTime)
	fmt.Printf("Time taken for query hospitalID (%s) policyID (%s): %s\n", hospitalID, policyID, duration)

	return asset, nil
}


// retrieveFromRegionalBC retrieves asset data from the regional blockchain using the provided path and asset ID.
func retrieveFromRegionalBC(ctx contractapi.TransactionContextInterface, rccName string, policyID string) (*RegionalAsset, error) {
	queryArgs := [][]byte{[]byte("ReadAsset"), []byte(policyID)}
	response := ctx.GetStub().InvokeCh aincode(rccName, queryArgs, "mychannel")

	if response.GetStatus() != shim.OK {
		return nil, fmt.Errorf("failed to retrieve asset data from regional blockchain: %s", response.GetMessage())
	}

	var regionalAsset RegionalAsset
	err := json.Unmarshal(response.Payload, &regionalAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal regional asset data: %v", err)
	}

	fmt.Printf("[REG] PolicyID: %s, Owner: %s\n", regionalAsset.ID, regionalAsset.Owner)
	return &regionalAsset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, hospitalID string, rccName string) error {
	exists, err := s.AssetExists(ctx, hospitalID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", hospitalID)
	}

	// overwriting original asset with new asset
	asset := GlobalAsset{
		HospitalID:     hospitalID,
		RegionalCCName: rccName,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hospitalID, assetJSON)
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
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newrccName string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.RegionalCCName = newrccName
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*GlobalAsset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*GlobalAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset GlobalAsset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (t *SmartContract) QueryAssetsByPolicyAndHospital(ctx contractapi.TransactionContextInterface, policyID string, hospitalID string) (*RegionalAsset, error) {
	startTime := time.Now()
	queryString := fmt.Sprintf(`{"selector":{"policyID":"%s","hospitalID":"%s"}}`, policyID, hospitalID)
	globalAsset, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("query global failed: %s", err)
	}

	fmt.Printf("[GLO] PolicyID: %s, HospitalID: %s, RegionalCC: %s\n", policyID, globalAsset.HospitalID, globalAsset.RegionalCCName)
	actualAsset, err := retrieveFromRegionalBC(ctx, globalAsset.RegionalCCName, policyID)
	if err != nil {
		return nil, fmt.Errorf("query regional failed: %s", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Time taken for query for policyID (%s): %s\n", policyID, duration)

	return actualAsset, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) (*GlobalAsset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*GlobalAsset, error) {
	// var assets []*GlobalAsset
	// for resultsIterator.HasNext() {
	// 	queryResult, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	var asset GlobalAsset
	// 	err = json.Unmarshal(queryResult.Value, &asset)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	assets = append(assets, &asset)
	// }

	// return assets[0], nil

	queryResult, err := resultsIterator.Next()
	if err != nil {
		return nil, err
	}
	var asset GlobalAsset
	err = json.Unmarshal(queryResult.Value, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating globalCC chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting globalCC chaincode: %v", err)
	}
}
