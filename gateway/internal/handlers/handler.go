package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"gateway/internal/fabric"
)

// Define the path to the CSV file
const csvFilePath = "../gateway/internal/index/hospital_index_table.csv"

var (
	chaincodeMap     map[string]string
	chaincodeMapOnce sync.Once
)

// ReadPPHandler handles the /readPP/ endpoint
func ReadPPHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	hospitalID := r.URL.Query().Get("hospitalID")
	policyID := r.URL.Query().Get("policyID")

	// Load hospital chaincode mapping from CSV (if not already loaded)
	chaincodeName, err := getChaincodeName(hospitalID)
	if err != nil {
		http.Error(w, "Failed to get chaincode name: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Call the Fabric network to retrieve data
	result, err := fabric.QueryAsset(chaincodeName, policyID)
	if err != nil {
		http.Error(w, "Failed to query asset from Fabric network: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract metadata from the result and redirect
	// Assuming metadata is a URL to redirect to
	var policy struct {
		ID        string   `json:"ID"`
		Owner     string   `json:"owner"`
		AuthRoles []string `json:"authRoles"`
		Grant     string   `json:"grant"`
		Metadata  string   `json:"metadata"`
	}

	if err := json.Unmarshal(result, &policy); err != nil {
		http.Error(w, "Failed to parse metadata", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate that the response is HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Display metadata URL with label and clickable link
	fmt.Fprintf(w, "HospitalID: %s<br>ChaincodeName: %s<br>\n", hospitalID, chaincodeName)
	fmt.Fprintf(w, "Data URL: <a href=\"%s\">%s</a><br><br>\n", policy.Metadata, policy.Metadata)

	// Return json
	json.NewEncoder(w).Encode(policy)
}

// getChaincodeName returns the chaincode name for the given hospital ID from the CSV file
func getChaincodeName(hospitalID string) (string, error) {
	// Load hospital chaincode mapping from CSV (if not already loaded)
	chaincodeMapOnce.Do(func() {
		chaincodeMap = make(map[string]string)
		file, err := os.Open(csvFilePath)
		if err != nil {
			fmt.Println("Failed to open CSV file:", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Failed to read CSV file:", err)
			return
		}

		for _, line := range lines {
			if len(line) >= 2 {
				chaincodeMap[line[0]] = line[1]
			}
		}
		fmt.Println("Chaincode map loaded successfully.")
	})

	// Retrieve chaincode name from the cached map
	chaincodeName, ok := chaincodeMap[hospitalID]
	if !ok {
		return "", errors.New("\nchaincode name not found for hospital ID: " + hospitalID)
	}
	return chaincodeName, nil
}

func PrintChaincodeMap() {
	fmt.Println("Chaincode Map:")
	for hospitalID, chaincodeName := range chaincodeMap {
		fmt.Printf("Hospital ID: %s, Chaincode Name: %s\n", hospitalID, chaincodeName)
	}
}
