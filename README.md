# cross-region-bc

## Linux Installation

Fabric Installation
```bash
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

./install-fabric.sh d s b
```

Copy Binary Files 
```bash
cp -r ~/fabric-samples/bin ~/cross-region-bc/
```

Go Installation
```bash
sudo apt update 
sudo apt install golang-go

cd ../crosschain/regional1
nano go.mod # change go version to match your machine
GO111MODULE=on go mod vendor
```

Time Library (gdate)
- In `indexer_execute.sh` and `onebc_execute.sh`, change from `gdate` to `date`
```bash
sudo apt-get install coreutils
```

Permission Issue
```bash
sudo chown -R {user} ${PWD}/cross-region-bc/test-network/organizations
```

Setting Path 
```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```