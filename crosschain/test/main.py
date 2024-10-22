import subprocess
import time
from concurrent.futures import ThreadPoolExecutor

# Function to execute a Fabric peer query and measure time for single blockchain
def query_single_chaincode(asset_id):
    try:
        start_time = time.time()
        query_command = f'peer chaincode query -C mychannel -n regionalCC1 -c \'{{"Args":["ReadAsset", "{asset_id}"]}}\''
        result = subprocess.run(query_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        query_time = (time.time() - start_time) / 60  # Convert to milliseconds

        if result.returncode == 0 and result.stdout.strip():
            return {"status": "true", "time": query_time}
        else:
            print('Single BC query failed:', result.stderr)
            return {"status": "false", "time": query_time}

    except Exception as e:
        print('Single Blockchain Error:', e)
        query_time = (time.time() - start_time) / 60  # Convert to milliseconds
        return {"status": "error", "time": query_time, "error": str(e)}

# Function to perform load test for both single and multiple blockchain models
def run_load_test(num_queries, num_runs, warmup_runs=1, sleep_duration=3, max_workers=10):
    asset_ids = [f"pc{(i % 100000) + 1}" for i in range(num_queries)]  # Asset IDs to query

    # # Warm-up phase for single blockchain
    # for _ in range(warmup_runs):
    #     results = {}
    #     start_time = time.time()
    #     with ThreadPoolExecutor(max_workers=max_workers) as executor:
    #         futures = {executor.submit(query_single_chaincode, asset_id): asset_id for asset_id in asset_ids}
    #         for future in futures:
    #             asset_id = futures[future]
    #             results[asset_id] = future.result()
        
    #     total_time = (time.time() - start_time) / 60
    #     print(f"Single Blockchain - Warm-up {_+1} >>> Query Times: {total_time:.4f} ms")
    #     time.sleep(sleep_duration)  # Sleep after each run

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
        
        total_time = (time.time() - start_time) / 60  # Convert to milliseconds
        total_single_time += total_time

        print(f"Single Blockchain - Run {_+1} >>> Query Times: {total_time:.4f} minutes")
        time.sleep(sleep_duration)  # Sleep after each run
    
    avg_single_time = total_single_time / num_runs

    print(f"Average time for {num_queries} queries in single blockchain: {avg_single_time:.4f} minutes")

# Main test loop
query_sizes = [10000, 50000, 100000, 150000, 200000]  # Different sizes to test
num_runs = 5  # Number of runs to calculate average
for size in query_sizes:
    print(f"\n--- Testing {size} queries ---")
    run_load_test(size, num_runs)
