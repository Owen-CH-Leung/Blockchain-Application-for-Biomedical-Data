#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# import utils
. scripts/getEnvVar.sh
. scripts/configUpdate.sh
. scripts/utility.sh

# NOTE: this must be run in a CLI container since it requires jq and configtxlator 
function createAnchorPeerUpdate() {
  Info "Fetching channel config for channel $CHANNEL_NAME"
  fetchChannelConfig $ORG $CHANNEL_NAME ${CORE_PEER_LOCALMSPID}config.json

  Info "Generating anchor peer update transaction for Org${ORG} on channel $CHANNEL_NAME"

  if [ $ORG == "Wales" ]; then
    HOST="peer0.WalesHospital.com"
    PORT=7051
  elif [ $ORG == "Margaret" ]; then
    HOST="peer0.MargaretHospital.com"
    PORT=9051
  else
    Error "Org${ORG} unknown"
  fi

  # Modify the configuration to append the anchor peer 
  jq '.channel_group.groups.Application.groups.'${CORE_PEER_LOCALMSPID}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'$HOST'","port": '$PORT'}]},"version": "0"}}' ${CORE_PEER_LOCALMSPID}config.json > ${CORE_PEER_LOCALMSPID}modified_config.json
  
  # Compute a config update, based on the differences between 
  # {orgmsp}config.json and {orgmsp}modified_config.json, write
  # it as a transaction to {orgmsp}anchors.tx
  createConfigUpdate ${CHANNEL_NAME} ${CORE_PEER_LOCALMSPID}config.json ${CORE_PEER_LOCALMSPID}modified_config.json ${CORE_PEER_LOCALMSPID}anchors.tx
  Info "Completed creation of AnchorPeer Update"
}

function updateAnchorPeer() {
  peer channel update -o orderer.HA.com:7050 --ordererTLSHostnameOverride orderer.HA.com -c $CHANNEL_NAME -f ${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile "$ORDERER_CA" >&log.txt
  res=$?
  errorCatch $res "Anchor peer update failed"
  Info "Updated AnchorPeer"
}

ORG=$1
CHANNEL_NAME=$2
if [ $ORG == "Wales" ]; then
  export CORE_PEER_ADDRESS=peer0.WalesHospital.com:7051
  setEnvironmentVariables $ORG
elif [ $ORG == "Margaret" ]; then
  export CORE_PEER_ADDRESS=peer0.MargaretHospital.com:9051
  setEnvironmentVariables $ORG
else
  Error "Unknown Organization. Export CORE_PEER_ADDRESS as empty string"
  export CORE_PEER_ADDRESS=""
fi

createAnchorPeerUpdate 
updateAnchorPeer 