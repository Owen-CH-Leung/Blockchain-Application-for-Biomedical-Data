#!/bin/bash
# createOrg 
#   - cryptogen (use crypto-config.yaml)
#   - docker compose fabric-ca (use docker-compose-ca.yaml, which has dependency on fabric-ca-order.yaml etc)
#   - enroll each user of each org
#   - create connection-org1.yaml (ccp-generate.sh)
# NetworkUp
#   - docker-compose to build the peer (use docker-compose-test-net.yaml) (need to define $DOCKER_SOCK)
# CreateChannel
#  - CreateGenesisBlock
#  - CreateChannel
#  - joinChannel (require to define $BLOCK_FILE)
#  - setAnchorPeer (which use scripts/setAnchorPeer.sh)
# DeployCC
# - run script.deployCC.sh
# Network Down

#docker-compose -f docker-ca.yaml up -d 2>&1
#upon compose, more files are generated

. utility.sh
. scripts/getEnvVar.sh

export CHANNEL_NAME=hospitalhk

function createFabricKeyMaterials(){
    Info "Creating Fabric Keys for HA, Wales Hospital & Margaret"
    cryptogen generate --config=./crypto-gen/crypto-config-demo.yaml --output="crypto-config"
    res=$?
    if [ $res -ne 0 ]; then
        Error "Failed to create fabric keys..."
    else
        Info "Successfully created fabric keys"
    fi
}

#Generate CCP script
function one_line_pem() {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function yaml_ccp() {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s/\${ORG_IDX}/$6/" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

function generateConectionYaml() {
    ORG="Wales"
    ORG_IDX=1
    P0PORT=7051
    CAPORT=7054
    PEERPEM=crypto-config/peerOrganizations/WalesHospital.com/tlsca/tlsca.WalesHospital.com-cert.pem
    CAPEM=crypto-config/peerOrganizations/WalesHospital.com/ca/ca.WalesHospital.com-cert.pem

    echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $ORG_IDX)" > crypto-config/peerOrganizations/WalesHospital.com/connection-org1.yaml

    ORG="Margaret"
    ORG_IDX=2
    P0PORT=9051
    CAPORT=8054
    PEERPEM=crypto-config/peerOrganizations/MargaretHospital.com/tlsca/tlsca.MargaretHospital.com-cert.pem
    CAPEM=crypto-config/peerOrganizations/MargaretHospital.com/ca/ca.MargaretHospital.com-cert.pem

    echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $ORG_IDX)" > crypto-config/peerOrganizations/MargaretHospital.com/connection-org2.yaml  
}

#Channel Creation
function createPeerOrderer(){
    Info "Started to build up peer & orderer"
    docker-compose -f ./peer/docker-compose-base-demo.yaml up -d 2>&1
    docker ps -a
    if [ $? -ne 0 ]; then
        Error "Unable to start network"
    else
        Info "Successfully created peer & orderer"
    fi
}

function generateGenesisBlock(){
    export FABRIC_CFG_PATH=${PWD}/configtx
    configtxgen -profile TwoHospitalsGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
    res=$?
    errorCatch $res "Failed to generate Genesis Block..."
    configtxgen -profile TwoHospitalsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
    res=$?
    errorCatch $res "Failed to generate ${CHANNEL_NAME}.tx..."
}

function createChannel(){
	setEnvironmentVariables Wales
    export BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"
    export FABRIC_CFG_PATH=${PWD}/config
	local rc=1
	local COUNTER=1
	while [ $rc -ne 0 -a $COUNTER -lt 5 ] ; do
		sleep 5
		peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.HA.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock $BLOCKFILE --tls --cafile $ORDERER_CA >&log.txt
        res=$?
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	errorCatch $res "Channel creation failed"
    if [ $res -eq 0 ]; then
        Info "Successfully created channel ${CHANNEL_NAME}"
    fi
}

function joinChannel(){
    export FABRIC_CFG_PATH=$PWD/config/
    setEnvironmentVariables $1
    export BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"
    local rc=1
    local COUNTER=1
    while [ $rc -ne 0 -a $COUNTER -lt 5 ] ; do
        sleep 5
        peer channel join -b $BLOCKFILE >&log.txt
        res=$?
        let rc=$res
        COUNTER=$(expr $COUNTER + 1)
    done
    errorCatch $res "After 5 attempts, peer0.$1 has failed to join channel '$CHANNEL_NAME' "
    if [ $res -eq 0 ]; then
        Info "peer0.$1 successfully join channel ${CHANNEL_NAME}"
    fi
}

function setAnchorPeer() {
    docker exec -it cli apk add ncurses
    docker exec -it cli ./scripts/setAnchorPeer.sh $1 $CHANNEL_NAME
}

function shutdownNetwork() {
    Info "Removing containers, volumes attached, and docker network..."
    docker kill $(docker ps -q) 
    docker rm $(docker ps -a -q)
    docker volume rm $(docker volume ls -q)
    docker network rm hospital-network
    sudo rm -rf crypto-config
    sudo rm -rf channel-artifacts
    sudo rm -rf system-genesis-block
    sudo rm hospitalCC.tar.gz
    rm log.txt
    Info "Completed"
}

function setupHyperledgerExplorer() {
    Info "Setting up Hyperledger Explorer..."
    cd ../hyperledger-explorer
    cp -r "../fabric-network/crypto-config" .
    docker-compose up -d
    docker ps -a | grep explorer
    if [ $? -ne 0 ]; then
        Error "Unable to start Hyperledger Explorer"
    else
        Info "Successfully created Hyperledger Explorer"
    fi
    cd ../fabric-network
}


if [[ $1 == "up" ]]; then 
    sudo chmod +x get_dependency.sh
    sudo chmod +x network.sh
    sudo chmod +x deployCC.sh
    sudo chmod +x utility.sh
    cd scripts

    sudo chmod +x configUpdate.sh
    sudo chmod +x getEnvVar.sh
    sudo chmod +x setAnchorPeer.sh
    sudo chmod +x utility.sh
    cd ..

    createFabricKeyMaterials
    generateGenesisBlock
    generateConectionYaml
    createPeerOrderer
    createChannel
    joinChannel Wales
    joinChannel Margaret
    setAnchorPeer Wales
    setAnchorPeer Margaret
    setupHyperledgerExplorer
fi 

if [[ $1 == "down" ]]; then 
    shutdownNetwork
fi
#Network Down
# 


#sudo apt install mlocate
#sudo apt install tree

#tree -H crypto-config > crypto-config.html
#rm -rf docker
#rm docker-20.10.7.tgz
#sudo rm -rf /usr/bin/docker-proxy
#sudo rm -rf /usr/bin/containerd
#sudo rm -rf /usr/bin/dockerd
#sudo rm -rf /usr/bin/docker
#sudo rm -rf /usr/bin/containerd-shim-runc-v2
#sudo rm -rf /usr/bin/ctr
#sudo rm -rf /usr/bin/docker-init
#sudo rm -rf /usr/bin/runc
#sudo rm -rf /usr/bin/containerd-shim
#sudo rm -rf /var/run/docker*
#sudo groupdel docker
#sudo rm -rf /etc/docker
#sudo rm -rf /var/lib/docker

# export FABRIC_CFG_PATH=${PWD}
# docker kill $(docker ps -q) 
# docker rm $(docker ps -a -q)
# docker volume rm $(docker volume ls -q)
# docker network rm hospital-network
# docker system prune -a
# ./network.sh up createChannel -c mychannel -ca

# rm -rf .fabric-ca-client
# . utility.sh

# cryptogen generate --config=crypto-config-demo.yaml --output="crypto-config-demo"

# configtxgen -profile TwoHospitalsGenesis -outputBlock ./channel-artifacts/mychannel.block -channelID mychannel

# Info "Removing previously-downloaded Hyperledger Fabric Images..."


# containerImage=$(docker ps -a | awk '{print $2}')
# containerID=($(docker ps -a | awk '{print $1}'))
# iter=1
# for imgID in $(docker images | grep hyperledger | awk '{print $3}')
# do
#     if [[ $containerImage =~ $imgID ]]; then
#         Alert "Found Image $imgID running. Container ID : ${containerID[${iter}]}"
#         Info "Killing Contaner..."
#         docker kill ${containerID[${iter}]}
#         Info "Removing Image $imgID..."
#         docker rmi $imgID
#         let iter=$iter+1
        
#     else
#         Info "Removing Image $imgID..."
#         docker rmi $imgID
#     fi
# done

# Info "Removing ALL Docker Volumes..."
# docker volume rm $(docker volume ls -q)
#shellcheck network.sh
