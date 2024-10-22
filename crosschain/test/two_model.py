import csv
import subprocess
import time
from concurrent.futures import ThreadPoolExecutor

# Load the hospital to chaincode mapping from CSV for multiple blockchain model
def load_chaincode_mapping(filename):
    hospital_chaincode_map = {}
    with open(filename, mode='r') as file:
        reader = csv.reader(file)
        for row in reader:
            hospital_chaincode_map[row[0]] = row[1]
    return hospital_chaincode_map

# Function to execute a Fabric peer query and measure time for single blockchain
def query_single_chaincode(asset_id):
    try:
        start_time = time.time()
        query_command = f'peer chaincode query -C mychannel -n regionalCC1 -c \'{{"Args":["ReadAsset", "{asset_id}"]}}\''
        result = subprocess.run(query_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        query_time = (time.time() - start_time) * 1000  # Convert to milliseconds

        if result.returncode == 0 and result.stdout.strip():
            return {"status": "true", "time": query_time}
        else:
            print('Single BC query failed:', result.stderr)
            return {"status": "false", "time": query_time}

    except Exception as e:
        print('Single Blockchain Error:', e)
        query_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        return {"status": "error", "time": query_time, "error": str(e)}

# Similarly update the query_multiple_chaincode function:
def query_multiple_chaincode(hospitalID, policyID, chaincode_map):
    try:
        start_time = time.time()
        chaincodeName = chaincode_map.get(hospitalID, None)

        if chaincodeName is None:
            query_time = (time.time() - start_time) * 1000  # Convert to milliseconds
            return {"status": "invalid", "time": query_time}

        query_command = f'peer chaincode query -C mychannel -n {chaincodeName} -c \'{{"Args":["ReadAsset", "pc{policyID}"]}}\''
        result = subprocess.run(query_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        query_time = (time.time() - start_time) * 1000  # Convert to milliseconds

        if result.returncode != 0 or "Error" in result.stderr or result.stdout.strip() == "":
            print('Multiple BC query failed:', result.stderr)
            return {"status": "false", "time": query_time}
        else:
            return {"status": "true", "time": query_time}

    except Exception as e:
        print('Multiple Blockchain Error:', e)
        query_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        return {"status": "error", "time": query_time, "error": str(e)}


# Function to perform load test for both single and multiple blockchain models
def run_load_test(num_queries, num_runs, warmup_runs=1, sleep_duration=3, max_workers=10):
    asset_ids = [f"pc{i+1}" for i in range(num_queries)]  # Asset IDs to query
    hospital_ids = [f"HP{i%8+1}" for i in range(num_queries)]  # Hospital IDs for multiple blockchain

    # Load the chaincode mapping for the multiple blockchain test
    chaincode_mapping = load_chaincode_mapping('hospital_chaincode_mapping.csv')

    # Warm-up phase for single blockchain
    for _ in range(warmup_runs):
        results = {}
        start_time = time.time()
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {executor.submit(query_single_chaincode, asset_id): asset_id for asset_id in asset_ids}
            for future in futures:
                asset_id = futures[future]
                results[asset_id] = future.result()
        
        total_time = (time.time() - start_time) * 1000
        print(f"Single Blockchain - Warm-up {_+1} >>> Query Times: {total_time:.4f} ms")
        time.sleep(sleep_duration)  # Sleep after each run

    # Test single blockchain
    total_single_time = 0
    for _ in range(num_runs):
        results = {}
        start_time = time.time()
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {executor.submit(query_single_chaincode, asset_id): asset_id for asset_id in asset_ids}
            for future in futures:
                asset_id = futures[future]
                results[asset_id] = future.result()
        
        total_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        total_single_time += total_time

        print(f"Single Blockchain - Run {_+1} >>> Query Times: {total_time:.4f} ms")
        time.sleep(sleep_duration)  # Sleep after each run
    
    avg_single_time = total_single_time / num_runs

    # Warm-up phase for multiple blockchain
    for _ in range(warmup_runs):
        results = {}
        start_time = time.time()
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {executor.submit(query_multiple_chaincode, hospital_ids[i], str(i+1), chaincode_mapping): str(i+1) for i in range(num_queries)}
            for future in futures:
                policyID = futures[future]
                results[policyID] = future.result()
        
        total_time = (time.time() - start_time) * 1000
        print(f"Multiple Blockchain - Warm-up {_+1} >>> Query Times: {total_time:.4f} ms")
        time.sleep(sleep_duration)  # Sleep after each run

    # Test multiple blockchains
    total_multi_time = 0
    for _ in range(num_runs):
        results = {}
        start_time = time.time()
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {executor.submit(query_multiple_chaincode, hospital_ids[i], str(i + 1), chaincode_mapping): str(i + 1) for i in range(num_queries)}
            for future in futures:
                policyID = futures[future]
                results[policyID] = future.result()
        
        total_time = (time.time() - start_time) * 1000  # Convert to milliseconds
        total_multi_time += total_time

        print(f"Multiple Blockchain - Run {_+1} >>> Query Times: {total_time:.4f} ms")
        time.sleep(sleep_duration)  # Sleep after each run
    
    avg_multi_time = total_multi_time / num_runs

    print(f"Average time for {num_queries} queries in single blockchain: {avg_single_time:.4f} ms")
    print(f"Average time for {num_queries} queries in multiple blockchain: {avg_multi_time:.4f} ms")

# Main test loop
query_sizes = [10, 100, 1000, 10000]  # Different sizes to test
num_runs = 5  # Number of runs to calculate average
for size in query_sizes:
    print(f"\n--- Testing {size} queries ---")
    run_load_test(size, num_runs)
