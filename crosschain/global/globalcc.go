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
	PolicyID       string `json:"policyID"`
	HospitalID     string `json:"hospitalID"`
	ID             string `json:"ID"`
	Owner          string `json:"owner"`
	FileCount      int    `json:"fileCount"`
	RegionalCCName string `json:"regionalCCName"`
}

type RegionalAsset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

func generateGlobalAssets(numAssets int) []GlobalAsset {
	var assets []GlobalAsset
	for i := 1; i <= numAssets; i++ {
		asset := GlobalAsset{
			ID:             fmt.Sprintf("glob%d", i),
			Owner:          "Dumb",
			HospitalID:     "HP0", // Assuming a default value for HospitalID
			PolicyID:       "PolicyID",
			FileCount:      0,                   // Assuming a default value for FileCount
			RegionalCCName: "anotherRegionalCC", // Assuming a default value for RegionalBCPath
		}
		if i%2 == 0 {
			// for test
			asset.Owner = "PATIENT!"
			asset.HospitalID = "HP1"
			asset.PolicyID = fmt.Sprintf("pc%d", i/2)
			asset.RegionalCCName = "regionalCC"
		}
		assets = append(assets, asset)
	}
	fmt.Printf("Successfully generated %d records\n", numAssets)
	return assets
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, numRows int) error {

	assets := generateGlobalAssets(numRows)

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
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, hospitalID string, policyID string, fileCount int, rccName string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := GlobalAsset{
		ID:             id,
		Owner:          owner,
		HospitalID:     hospitalID,
		PolicyID:       policyID,
		FileCount:      fileCount,
		RegionalCCName: rccName,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
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

// retrieveFromRegionalBC retrieves asset data from the regional blockchain using the provided path and asset ID.
func retrieveFromRegionalBC(ctx contractapi.TransactionContextInterface, rccName string, policyID string) (*RegionalAsset, error) {
	queryArgs := [][]byte{[]byte("ReadAsset"), []byte(policyID)}
	response := ctx.GetStub().InvokeChaincode(rccName, queryArgs, "mychannel")

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
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, hospitalID string, policyID string, fileCount int, rccName string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := GlobalAsset{
		ID:             id,
		Owner:          owner,
		HospitalID:     hospitalID,
		PolicyID:       policyID,
		FileCount:      fileCount,
		RegionalCCName: rccName,
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

	fmt.Printf("[GLO] GlobalID: %s, Owner: %s, PolicyID: %s, HospitalID: %s\n", globalAsset.ID, globalAsset.Owner, globalAsset.PolicyID, globalAsset.HospitalID)
	actualAsset, err := retrieveFromRegionalBC(ctx, globalAsset.RegionalCCName, globalAsset.PolicyID)
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
