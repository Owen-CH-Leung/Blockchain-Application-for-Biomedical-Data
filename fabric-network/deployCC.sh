#!/bin/bash

# source scripts/utils.sh

# CHANNEL_NAME=${1:-"mychannel"}
# CC_NAME=${2}
# CC_SRC_PATH=${3}
# CC_SRC_LANGUAGE=${4}
# CC_VERSION=${5:-"1.0"}
# CC_SEQUENCE=${6:-"1"}
# CC_INIT_FCN=${7:-"NA"}
# CC_END_POLICY=${8:-"NA"}
# CC_COLL_CONFIG=${9:-"NA"}
# DELAY=${10:-"3"}
# MAX_RETRY=${11:-"5"}
# VERBOSE=${12:-"false"}


# FABRIC_CFG_PATH=$PWD/../config/

# #User has not provided a name
# if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
#   fatalln "No chaincode name was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

# # User has not provided a path
# elif [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
#   fatalln "No chaincode path was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

# # User has not provided a language
# elif [ -z "$CC_SRC_LANGUAGE" ] || [ "$CC_SRC_LANGUAGE" = "NA" ]; then
#   fatalln "No chaincode language was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

# ## Make sure that the path to the chaincode exists
# elif [ ! -d "$CC_SRC_PATH" ]; then
#   fatalln "Path to chaincode does not exist. Please provide different path."
# fi

# CC_SRC_LANGUAGE=$(echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:])

# # do some language specific preparation to the chaincode before packaging
# if [ "$CC_SRC_LANGUAGE" = "go" ]; then
#   CC_RUNTIME_LANGUAGE=golang

#   infoln "Vendoring Go dependencies at $CC_SRC_PATH"
#   pushd $CC_SRC_PATH
#   GO111MODULE=on go mod vendor
#   popd
#   successln "Finished vendoring Go dependencies"

# else
#   fatalln "The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script. Supported chaincode languages are: go, java, javascript, and typescript"
#   exit 1
# fi

# INIT_REQUIRED="--init-required"
# # check if the init fcn should be called
# if [ "$CC_INIT_FCN" = "NA" ]; then
#   INIT_REQUIRED=""
# fi

# if [ "$CC_END_POLICY" = "NA" ]; then
#   CC_END_POLICY=""
# else
#   CC_END_POLICY="--signature-policy $CC_END_POLICY"
# fi

# if [ "$CC_COLL_CONFIG" = "NA" ]; then
#   CC_COLL_CONFIG=""
# else
#   CC_COLL_CONFIG="--collections-config $CC_COLL_CONFIG"
# fi

# import utils
. scripts/getEnvVar.sh
. scripts/utility.sh

CHANNEL_NAME=${1}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_VERSION=${4:-"1"}
CC_SEQUENCE=${5:-"1"}
CC_INIT_FCN=${6:-"NA"}
CC_END_POLICY=${7:-"NA"}
CC_COLL_CONFIG=${8:-"NA"}
INIT_REQUIRED=${9:-""}

export FABRIC_CFG_PATH=$PWD/config/
export CHANNEL_NAME=hospitalhk
export CC_NAME=hospitalCC
export CC_SRC_PATH=../fabric-chaincode/

function installCC() {
  ORG=$1
  setEnvironmentVariables ${ORG}
  peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --lang golang --label ${CC_NAME}_${CC_VERSION} >&cclog.txt
  res=$?
  errorCatch $res "Chaincode packaging has failed"
  if [ $res -eq 0 ]; then
      Info "Successfully packaged ${CC_NAME} as ${CC_NAME}.tar.gz"
  fi

  peer lifecycle chaincode install ${CC_NAME}.tar.gz >&cclog.txt
  res=$?
  errorCatch $res "Chaincode installation on peer0.${ORG} has failed"
  if [ $res -eq 0 ]; then
      Info "Successfully installed ${CC_NAME} on peer0.${ORG}"
  fi

  peer lifecycle chaincode queryinstalled >&cclog.txt
  if [ $res -eq 0 ]; then
      Info "Successfully query the installed CC on peer0.${ORG}"
  else
      Error "Query installed on peer0.${ORG} has failed"
  fi
  PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" cclog.txt)
  Info "Chaincode is packaged & installed. Queryinstallation is completed"
}

# approveForMyOrg VERSION PEER ORG
function makeapproval() {
  ORG=$1
  setEnvironmentVariables ${ORG}
  peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.HA.com --tls --cafile "$ORDERER_CA" --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&cclog.txt
  if [ $res -eq 0 ]; then
      Info "Chaincode definition approved on peer0.${ORG} on '$CHANNEL_NAME'"
  else
      Error "Chaincode definition NOT approved on peer0.${ORG} on '$CHANNEL_NAME'"
  fi
}

# checkCommitReadiness VERSION PEER ORG
function iscommitready() {
  ORG=$1
  shift 1
  setEnvironmentVariables $ORG
  Info "Checking the commit readiness of the chaincode definition on peer0.${ORG} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt 5 ]; do
    sleep 3
    Info "Attempting to check the commit readiness, Retry after 3 seconds."
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} --output json >&cclog.txt
    res=$?
    let rc=0
    for var in "$@"; do
      grep "$var" cclog.txt &>/dev/null || let rc=1
    done
    COUNTER=$(expr $COUNTER + 1)
  done
  if test $rc -eq 0; then
    Info "Chaincode definition is ready for commit on peer0.${ORG} on channel '$CHANNEL_NAME'"
  else
    Error "After 5 attempts, Chaincode definition is NOT ready for commit on peer0.${ORG}!"
  fi
}

# commitChaincodeDefinition VERSION PEER ORG (PEER ORG)...
function commitCC() {
  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  ORG1=$1
  ORG2=$2
  setEnvironmentVariables $ORG1
  tlscert1=${PEER0_ORG1_CA}
  tlscert2=${PEER0_ORG2_CA}
  lc1=${CORE_PEER_ADDRESS}
  setEnvironmentVariables $ORG2
  lc2=${CORE_PEER_ADDRESS}
  
  peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.HA.com --tls --cafile "$ORDERER_CA" --channelID $CHANNEL_NAME --name ${CC_NAME} --peerAddresses ${lc1} --tlsRootCertFiles "${tlscert1}" --peerAddresses ${lc2} --tlsRootCertFiles "${tlscert2}" --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&cclog.txt
  res=$?
  errorCatch $res "Chaincode definition commit failed on peer0.${ORG} on channel '$CHANNEL_NAME' failed"
  Info "Chaincode definition committed on channel '$CHANNEL_NAME'"
}
# peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/ubuntu/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem 
#--channelID mychannel --name basic --peerAddresses localhost:7051 
#--tlsRootCertFiles 
#/home/ubuntu/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
#--peerAddresses localhost:9051 
#--tlsRootCertFiles /home/ubuntu/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
#--version 1.0 --sequence 1

# queryCommitted ORG
function queryCommitted() {
  ORG=$1
  setEnvironmentVariables $ORG
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  Info "Querying chaincode definition on peer0.${ORG} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt 5 ]; do
    sleep 3
    Info "Attempting to Query committed status on peer0.${ORG}, Retry after 3 seconds."
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&cclog.txt
    test $res -eq 0 && VALUE=$(cat cclog.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat cclog.txt
  if test $rc -eq 0; then
    Info "Query chaincode definition successful on peer0.${ORG} on channel '$CHANNEL_NAME'"
  else
    Error "After 5 attempts, Query chaincode definition result on peer0.${ORG} is INVALID!"
  fi
}

installCC Wales
installCC Margaret
makeapproval Wales
iscommitready "Wales" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": false"
iscommitready "Margaret" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": false"

makeapproval Margaret
iscommitready "Wales" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": true"
iscommitready "Margaret" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": true"

commitCC Wales Margaret
queryCommitted Wales
queryCommitted Margaret

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.HA.com --tls --cafile ${PWD}/crypto-config/ordererOrganizations/HA.com/orderers/orderer.HA.com/msp/tlscacerts/tlsca.HA.com-cert.pem -C hospitalhk -n hospitalCC --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/WalesHospital.com/peers/peer0.WalesHospital.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/MargaretHospital.com/peers/peer0.MargaretHospital.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

# peer chaincode query -C hospitalhk -n hospitalCC -c '{"Args":["GetAllAssets"]}'
## package the chaincode

## Install chaincode on peer0.org1 and peer0.org2
# Info "Installing chaincode on peer0.WalesHospital..."
# installCC "Wales"
# Info "Install chaincode on peer0.MargaretHospital..."
# installCC "Margaret"

# ## approve the definition for org1
# makeapproval "Wales"

# ## check whether the chaincode definition is ready to be committed
# ## expect org1 to have approved and org2 not to
# iscommitready "Wales" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": false"
# iscommitready "Margaret" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": false"

# ## now approve also for org2
# makeapproval "Margaret"

# ## check whether the chaincode definition is ready to be committed
# ## expect them both to have approved
# iscommitready "Wales" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": true"
# iscommitready "Margaret" "\"WalesHospitalMSP\": true" "\"MargaretHospitalMSP\": true"

# ## now that we know for sure both orgs have approved, commit the definition
# commitCC "Wales" "Margaret"

# ## query on both orgs to see that the definition committed successfully
# queryCommitted "Wales"
# queryCommitted "Margaret"

# ## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
# ## method defined
# if [ "$CC_INIT_FCN" = "NA" ]; then
#   infoln "Chaincode initialization is not required"
# else
#   chaincodeInvokeInit 1 2
# fi

# exit 0