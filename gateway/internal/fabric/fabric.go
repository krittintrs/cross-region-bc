package fabric

import (
	"fmt"
	"os/exec"
	"os"
)

// QueryAsset queries the asset from the Fabric network
func QueryAsset(hospitalID, policyID string) ([]byte, error) {
	// Change directory to the desired location
	if err := os.Chdir("../test-network"); err != nil {
		return nil, fmt.Errorf("failed to change directory: %v", err)
	}

	// Construct the command to execute
	cmd := fmt.Sprintf("peer chaincode query -C mychannel -n %s -c '{\"Args\":[\"ReadAsset\", \"%s\"]}'", hospitalID, policyID)

	// Execute the command
	cmdOutput, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s\nERROR: %v\nCommand output: %s", cmd, err, cmdOutput)
	}

	return cmdOutput, nil
}
