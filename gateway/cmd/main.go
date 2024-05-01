package main

import (
	"fmt"
	"log"
	"os"

	"gateway/internal/server"
)

func main() {
    // Initial Setup
    Setup(1)

	// Initialize and start the HTTP server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func Setup(orgID int) {
    // Change Directory
    os.Chdir("../test-network")
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current working directory:", err)
        return
    }

    // Modify the PWD env var
    os.Setenv("PWD", cwd)

    // Modify the PATH env var
    newPath := fmt.Sprintf("%s/../bin:%s", cwd, os.Getenv("PATH"))
    os.Setenv("PATH", newPath)

    // Set peer cmd path
	SetPeerCmdPath()

	// Set the path for Org1
	SetPeerPath(orgID)
}

func SetPeerCmdPath() {
	// Set the required environment variables
    os.Setenv("PATH", os.Getenv("PWD")+"/../bin:"+os.Getenv("PATH"))
    os.Setenv("FABRIC_CFG_PATH", os.Getenv("PWD")+"/../config/")

	fmt.Println("Peer chaincode cmd path set")
}

// SetPeerPath sets the peer path based on the organization identifier.
func SetPeerPath(orgID int) {
    var mspID, peerAddress, tlsCertPath string

    switch orgID {
    case 1:
        mspID = "Org1MSP"
        peerAddress = "localhost:7051"
        tlsCertPath = os.Getenv("PWD") + "/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
    case 2:
        mspID = "Org2MSP"
        peerAddress = "localhost:9051"
        tlsCertPath = os.Getenv("PWD") + "/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
    default:
        log.Fatalf("Invalid organization ID: %d", orgID)
    }

    // Set environment variables
    os.Setenv("CORE_PEER_TLS_ENABLED", "true")
    os.Setenv("CORE_PEER_LOCALMSPID", mspID)
    os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", tlsCertPath)
    os.Setenv("CORE_PEER_MSPCONFIGPATH", os.Getenv("PWD")+"/organizations/peerOrganizations/org"+fmt.Sprintf("%d", orgID)+".example.com/users/Admin@org"+fmt.Sprintf("%d", orgID)+".example.com/msp")
    os.Setenv("CORE_PEER_ADDRESS", peerAddress)

    fmt.Printf("Peer path set for Org%d\n", orgID)
}
