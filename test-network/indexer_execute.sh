#!/bin/bash

# Capture the start time in nanoseconds
startTime=$(gdate +%s%N)

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <hospitalID> <policyID>"
    exit 1
fi

hospitalID=$1
policyID=$2

# Look up the chaincode name based on the hospitalID
chaincodeName=$(awk -F',' -v id="$hospitalID" '$1==id {print $2}' hospital_chaincode_mapping.csv)

# Check if the hospitalID is valid
if [ -z "$chaincodeName" ]; then
    echo "Invalid hospitalID: $hospitalID"
    
    # Capture the end time in nanoseconds
    endTime=$(gdate +%s%N)

    # Calculate the time difference in milliseconds
    duration=$(( ($endTime - $startTime)/1000000 ))

    # Print the total time taken
    echo "Total time taken: $duration milliseconds" 
    
    exit 1
fi

# Construct the query command
queryCommand="peer chaincode query -C mychannel -n $chaincodeName -c '{\"Args\":[\"ReadAsset\", \"pc$policyID\"]}'"

# Execute the query command
echo "Executing command: $queryCommand"
response=$(eval $queryCommand)

# Capture the end time in nanoseconds
endTime=$(gdate +%s%N)

# Calculate the time difference in milliseconds
duration=$(( ($endTime - $startTime)/1000000 ))

# Print the response
echo "$response"

# Print the total time taken
echo "Total time taken: $duration milliseconds"
