#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <policyID>"
    exit 1
fi

# Set variables for policyID
policyID=$1

# Capture the start time in nanoseconds
startTime=$(gdate +%s%N)

# Construct the query command
queryCommand="peer chaincode query -C mychannel -n regionalCC1 -c '{\"Args\":[\"ReadAsset\", \"pc$policyID\"]}'"

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
