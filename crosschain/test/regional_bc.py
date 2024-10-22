import csv
import subprocess
import threading
import time

# Load the hospital to chaincode mapping from CSV
def load_chaincode_mapping(filename):
    hospital_chaincode_map = {}
    with open(filename, mode='r') as file:
        reader = csv.reader(file)
        for row in reader:
            hospital_chaincode_map[row[0]] = row[1]
    return hospital_chaincode_map

# Function to execute a Fabric peer query and measure time
def query_chaincode(hospitalID, policyID, chaincode_map, results):
    try:
        # Start the timer
        start_time = time.time()
        
        # Look up the chaincode name based on the hospitalID
        chaincodeName = chaincode_map.get(hospitalID, None)

        # Check if the hospitalID is valid
        if chaincodeName is None:
            end_time = time.time()
            query_time = (end_time - start_time) * 1000  # Convert to milliseconds
            results[policyID] = {"status": "invalid", "time": query_time}
            return

        # Construct the query command
        query_command = f'peer chaincode query -C mychannel -n {chaincodeName} -c \'{{"Args":["ReadAsset", "pc{policyID}"]}}\''
        
        # Execute the command and capture the output
        result = subprocess.run(query_command, shell=True, capture_output=True, text=True)
        
        # End the timer
        end_time = time.time()
        query_time = (end_time - start_time) * 1000  # Convert to milliseconds
        
        # Check if the response contains an error
        if result.returncode != 0 or "Error" in result.stderr or result.stdout.strip() == "":
            results[policyID] = {"status": "false", "time": query_time}
        else:
            results[policyID] = {"status": "true", "time": query_time}
    
    except Exception as e:
        end_time = time.time()
        query_time = (end_time - start_time) * 1000  # Convert to milliseconds
        results[policyID] = {"status": "error", "time": query_time, "error": str(e)}

# List of hospital and policy IDs for queries
queries = [("HP1", "1"), ("HP2", "2"), ("HP3", "3"), ("HP4", "4"), ("HP5", "5")]

# Load the chaincode mapping from the CSV file
chaincode_mapping = load_chaincode_mapping('hospital_chaincode_mapping.csv')

# Dictionary to store results and times for each query
results = {}

# Create a thread for each query
threads = []
start_time_all = time.time()  # Start timer for the entire process
for hospitalID, policyID in queries:
    thread = threading.Thread(target=query_chaincode, args=(hospitalID, policyID, chaincode_mapping, results))
    threads.append(thread)
    thread.start()

# Wait for all threads to complete
for thread in threads:
    thread.join()

end_time_all = time.time()  # End timer for the entire process
total_time = (end_time_all - start_time_all) * 1000  # Convert to milliseconds

# Print the results and time taken for each query
for policyID in results:
    result = results[policyID]
    print(f"Policy ID: {policyID}\tStatus: {result['status']}\tTime: {result['time']:.4f} ms")

# Print the total time taken
print(f"\nTotal time for all queries: {total_time:.4f} ms")
