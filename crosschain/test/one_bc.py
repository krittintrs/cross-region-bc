import subprocess
import threading
import time

# Function to execute a Fabric peer query and measure time
def query_chaincode(asset_id, results):
    try:
        # Start the timer
        start_time = time.time()
        
        # Construct the query command
        query_command = f'peer chaincode query -C mychannel -n regionalCC1 -c \'{{"Args":["ReadAsset", "{asset_id}"]}}\''
        
        # Execute the command
        result = subprocess.run(query_command, shell=True, capture_output=True, text=True)
        
        # End the timer
        end_time = time.time()
        query_time = (end_time - start_time) * 1000  # Convert to milliseconds
        
        # Store the result status
        if result.returncode == 0 and result.stdout.strip():
            status = "true"
        else:
            status = "false"
        
        # Save results and query time
        results[asset_id] = {"status": status, "time": query_time}
    
    except Exception as e:
        query_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        results[asset_id] = {"status": "error", "time": query_time, "error": str(e)}

# List of asset IDs to query concurrently
asset_ids = ["pc1", "pc2", "pc3", "pc4", "pc5"]

# Dictionary to store results and times for each query
results = {}

# Create a thread for each query
threads = []
start_time_all = time.time()  # Start timer for the entire process
for asset_id in asset_ids:
    thread = threading.Thread(target=query_chaincode, args=(asset_id, results))
    threads.append(thread)
    thread.start()

# Wait for all threads to complete
for thread in threads:
    thread.join()

end_time_all = time.time()  # End timer for the entire process
total_time = (end_time_all - start_time_all) * 1000  # Convert to milliseconds

# Print the results and time taken for each query in the same format as the regional script
for asset_id in asset_ids:
    print(f"Policy ID: {asset_id[-1]}\tStatus: {results[asset_id]['status']}\tTime: {results[asset_id]['time']:.4f} ms")

# Print the total time taken
print(f"\nTotal time for all queries: {total_time:.4f} ms")
