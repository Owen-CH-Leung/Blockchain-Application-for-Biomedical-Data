#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# import utils
#. scripts/getEnvVar.sh
. scripts/utility.sh
# fetchChannelConfig <org> <channel_id> <output_json>
# Writes the current channel config for a given channel to a JSON file
# NOTE: this must be run in a CLI container since it requires configtxlator 
function fetchChannelConfig() {
  ORG=$1
  CHANNEL=$2
  OUTPUT=$3

  setEnvironmentVariables $ORG

  Info "Fetching the most recent configuration block for the channel"
  peer channel fetch config config_block.pb -o orderer.HA.com:7050 --ordererTLSHostnameOverride orderer.HA.com -c $CHANNEL --tls --cafile "$ORDERER_CA"
  res=$?
  errorCatch $res "Fetching config from peer failed"

  Info "Decoding config block to JSON and isolating config to ${OUTPUT}"
  configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  res=$?
  errorCatch $res "Decoding config block failed"
}

# createConfigUpdate <channel_id> <original_config.json> <modified_config.json> <output.pb>
# Takes an original and modified config, and produces the config update tx
# which transitions between the two
# NOTE: this must be run in a CLI container since it requires configtxlator 
function createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  res=$?
  errorCatch $res "creating original_config.pb failed"

  configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  res=$?
  errorCatch $res "creating modified_config.pb failed"

  configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb
  res=$?
  errorCatch $res "creating config_update.pb failed"

  configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  res=$?
  errorCatch $res "creating config_update.json failed"

  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
  res=$?
  errorCatch $res "converting config_update_in_envelope.json failed"
}

# signConfigtxAsPeerOrg <org> <configtx.pb>
# Set the peerOrg admin of an org and sign the config update
signConfigtxAsPeerOrg() {
  ORG=$1
  CONFIGTXFILE=$2
  setEnvironmentVariables $ORG
  peer channel signconfigtx -f "${CONFIGTXFILE}"
}