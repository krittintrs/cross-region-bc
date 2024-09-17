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
response=$(eval $queryCommand 2>&1)  # Capture both stdout and stderr

# Capture the end time in nanoseconds
endTime=$(gdate +%s%N)

# Calculate the time difference in microseconds
duration=$(( ($endTime - $startTime)/1000 ))

# echo "$response"

# Check if the response contains an error message
if echo "$response" | grep -q "Error" || [ -z "$response" ]; then
    status="false"
else
    status="true"
fi

# Print the result in tab-separated format
echo -e "$policyID\t$status\t$duration"
