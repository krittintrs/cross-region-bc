#!/bin/bash

# Define the list of data sizes
data_sizes=("10" "100" "1000" "10000" "100000" "300000")

# Loop over each data size
for data in "${data_sizes[@]}"; do
    echo "Processing data size: $data"
    
    # Run the setup and test scripts with the current data size
    ./onebc_network_setup.sh
    echo "Network Setup >> DONE"

    ./onebc_data.sh "$data"
    echo "Insert Data >> DONE"

    sleep 5

    ./onebc_test.sh "$data" > "../result/onebc/data${data}.txt" 2>&1
    echo "Query Data >> DONE"

    ./network.sh down
    echo "[ ONE BC - $data data ] >> DONE"
done
